package ssr

import (
	"context"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
)

type SessionValidator interface {
	ValidateSession(ctx context.Context, token string) (string, error)
}

type contextKey string

const UserIDKey contextKey = "userID"

func LoggingMiddleware(baseLogger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLogger := baseLogger.With("method", r.Method, "path", r.URL.Path)
		ctx := logger.WithContext(r.Context(), reqLogger)
		reqWithContext := r.WithContext(ctx)

		next.ServeHTTP(w, reqWithContext)
	})
}

func RequireAuth(validator SessionValidator, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := validator.ValidateSession(r.Context(), cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "session_token",
				MaxAge: -1,
				Path:   "/",
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		reqWithCtx := r.WithContext(ctx)

		next.ServeHTTP(w, reqWithCtx)
	}
}
