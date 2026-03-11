package ssr

import (
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
)

func LoggingMiddleware(baseLogger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := baseLogger.With("method", r.Method, "path", r.URL.Path)
		ctx := logger.IntoContext(r.Context(), reqLogger)
		reqWithContext := r.WithContext(ctx)

		next.ServeHTTP(w, reqWithContext)
	})
}
