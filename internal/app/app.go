package app

import (
	"context"
	"fmt"
	"learningflow/internal/transport/ssr"
	"learningflow/pkg/config"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
)

type App struct {
	HTTPPort string
	Logs     *slog.Logger
}

func New(ctx context.Context) (*App, error) {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Error("app.New, failed to load config",
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("app.New, failed to load config: %w", err)
	}

	logger.Setup(cfg.AppEnv)

	logs := logger.With("serice", "learningflow")
	logger.IntoContext(ctx, logs)
	logs.Info("initializing layers", "env", cfg.AppEnv, "port", cfg.HTTPPort)

	return &App{HTTPPort: cfg.HTTPPort, Logs: logs}, nil
}

func (a *App) Run(ctx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", ssr.LoggingMiddleware(a.Logs, ssr.HandleHome))
	mux.HandleFunc("/category/{id}", ssr.LoggingMiddleware(a.Logs, ssr.HandleCategory))
	mux.HandleFunc("/tutors", ssr.LoggingMiddleware(a.Logs, ssr.HandleTutors))
	mux.HandleFunc("/tutor/{id}", ssr.LoggingMiddleware(a.Logs, ssr.HandleTutor))

	fileserver := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileserver))

	a.Logs.Info("App.Run starting server", "port", a.HTTPPort)
	if err := http.ListenAndServe(":"+a.HTTPPort, mux); err != nil {
		a.Logs.Error("App.Run failed to start server", "port", a.HTTPPort)
	}
}
