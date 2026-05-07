package app

import (
	"context"
	"fmt"
	"learningflow/internal/migrations"
	"learningflow/internal/repository/db"
	"learningflow/internal/service"
	"learningflow/internal/transport/ssr"
	"learningflow/internal/transport/ssr/auth"
	"learningflow/internal/transport/ssr/catalog"
	"learningflow/pkg/closer"
	"learningflow/pkg/config"
	"learningflow/pkg/logger"
	"learningflow/pkg/migrator"
	"learningflow/pkg/postgress"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	HTTPPort   string
	pool       *pgxpool.Pool
	closer     *closer.Closer
	logs       *slog.Logger
	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return nil, fmt.Errorf("app.New, failed to load config: %w", err)
	}

	logger.Setup(cfg.AppEnv)

	logs := logger.With("serice", "learningflow")
	logger.WithContext(ctx, logs)
	logs.Info("initializing layers", "env", cfg.AppEnv, "port", cfg.HTTPPort)
	ctx = logger.WithContext(ctx, logs)
	pool, err := postgress.NewPool(ctx, cfg.PGDsn)
	if err != nil {
		return nil, fmt.Errorf("app.New, failed to initialize postgres pool: %w", err)
	}
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()
	m, err := migrator.EmbedMigrations(sqlDB, migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("app.New, failed to initialize migrator: %w", err)
	}
	if err := m.Up(); err != nil {
		return nil, fmt.Errorf("app.New, failed to initialize migrations: %w", err)
	}

	userRepo := db.NewUserRepository(pool)
	sessionRepo := db.NewSessionRepository(pool)
	tutorRepo := db.NewTutorRepo(pool)
	subjectRepo := db.NewSubjectRepo(pool)

	authSvc := service.NewAuthService(userRepo, sessionRepo)
	tutorSvc := service.NewTutorService(tutorRepo)
	subjectSvc := service.NewSubjectService(subjectRepo)

	catalogHandler := catalog.NewHandlerCatalog(subjectSvc, tutorSvc)
	authHandler := auth.NewAuthHandler(authSvc)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", ssr.LoggingMiddleware(logs, catalogHandler.HandleHome))
	mux.HandleFunc("GET /category/{id}", ssr.LoggingMiddleware(logs, catalogHandler.HandleCategory))
	mux.HandleFunc("GET /tutors", ssr.LoggingMiddleware(logs, catalogHandler.HandleTutors))
	mux.HandleFunc("GET /tutor/{id}", ssr.LoggingMiddleware(logs, catalogHandler.HandleTutor))

	mux.HandleFunc("GET /login", ssr.LoggingMiddleware(logs, authHandler.ShowLogin))
	mux.HandleFunc("POST /login", ssr.LoggingMiddleware(logs, authHandler.ProcessLogin))
	mux.HandleFunc("GET /register", ssr.LoggingMiddleware(logs, authHandler.ShowRegister))
	mux.HandleFunc("POST /register", ssr.LoggingMiddleware(logs, authHandler.ProcessRegister))

	fileserver := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileserver))

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	cl := closer.New()

	cl.Add(func(ctx context.Context) error {
		slog.Info("closing database connection pool")
		pool.Close()
		return nil
	})

	cl.Add(func(ctx context.Context) error {
		slog.Info("closing http server")
		return httpServer.Shutdown(ctx)
	})

	return &App{
		HTTPPort:   cfg.HTTPPort,
		pool:       pool,
		logs:       logs,
		httpServer: httpServer,
		closer:     cl,
	}, nil
}

func (a *App) Run() {
	errCh := make(chan error)

	go func() {
		a.logs.Info("starting http server")
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("app.Run: %w", err)
		}
	}()

	a.logs.Info("App.Run starting server", "port", a.HTTPPort)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		a.logs.Error("app.run server startup failed", "error", err)
	case sig := <-quit:
		a.logs.Info("app.run server shutdown", "signal", sig)
	}

	a.logs.Info("shutting down servers")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		a.logs.Error("app.Run shutdown failed", "error", err)
	}

	fmt.Println("Server Stopped")
}
