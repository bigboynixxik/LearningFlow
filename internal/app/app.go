package app

import (
	"context"
	"learningflow/internal/transport/ssr"
	"learningflow/pkg/config"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
)

type App struct {
	HTTPPort string
}

func New(ctx context.Context) *App {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Error("app.New, failed to load config",
			slog.String("error", err.Error()))
	}

	logger.Setup(cfg.AppEnv)

	logs := logger.With("serice", "learningflow")
	logger.IntoContext(ctx, logs)
	logger.FromContext(ctx).Info("initializing layers", "env", cfg.AppEnv, "port", cfg.HTTPPort)

	return &App{HTTPPort: cfg.HTTPPort}
}

func (a *App) Run(ctx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ssr.HandleHome)
	logger.FromContext(ctx).Info("App.Run starting server", "port", a.HTTPPort)
	if err := http.ListenAndServe(":"+a.HTTPPort, mux); err != nil {
		logger.FromContext(ctx).Error("App.Run failed to start server", "port", a.HTTPPort)
	}
}
