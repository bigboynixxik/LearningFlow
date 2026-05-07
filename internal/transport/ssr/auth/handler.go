package auth

import (
	"context"
	"errors"
	"html/template"
	"learningflow/internal/models"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
	"time"
)

type ServiceAuth interface {
	Register(ctx context.Context, email, plainPassword string, role models.Role) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
	ValidateSession(ctx context.Context, token string) (string, error)
}

type HandlerAuth struct {
	authService ServiceAuth
}

func NewAuthHandler(authService ServiceAuth) *HandlerAuth {
	return &HandlerAuth{authService: authService}
}

func (h *HandlerAuth) ShowLogin(w http.ResponseWriter, r *http.Request) {
	h.renderAuthForm(w, r, "login.html", "", http.StatusOK)
}

func (h *HandlerAuth) ProcessLogin(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		h.renderAuthForm(w, r, "login.html", "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, err := h.authService.Login(r.Context(), email, password)
	if err != nil {
		l.Error("auth.ProcessLogin authService login error", slog.String("err", err.Error()))
		h.renderAuthForm(w, r, "login.html", "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	h.setSessionCookie(w, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerAuth) ShowRegister(w http.ResponseWriter, r *http.Request) {
	h.renderAuthForm(w, r, "register.html", "", http.StatusOK)
}

func (h *HandlerAuth) ProcessRegister(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		l.Error("auth.ProcessRegister ParseForm error", slog.String("err", err.Error()))
		h.renderAuthForm(w, r, "register.html", "Ошибка обработки формы", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	if password != passwordConfirm {
		h.renderAuthForm(w, r, "register.html", "Пароли не совпадают", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Register(r.Context(), email, password, models.RoleStudent)
	if err != nil {
		if errors.Is(err, models.ErrAlreadyExists) {
			h.renderAuthForm(w, r, "register.html", "Пользователь с таким email уже существует", http.StatusConflict)
			return
		}

		l.Error("auth.ProcessRegister authService register error", slog.String("err", err.Error()))
		h.renderAuthForm(w, r, "register.html", "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	h.setSessionCookie(w, token)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *HandlerAuth) setSessionCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

type TemplateAuthData struct {
	Title string
	Error string
}

func (h *HandlerAuth) renderAuthForm(w http.ResponseWriter, r *http.Request, tmplName string, errMsg string, statusCode int) {
	l := logger.FromContext(r.Context())

	tmpl, err := template.ParseFiles("web/templates/" + tmplName)
	if err != nil {
		// Логируем ошибку парсинга
		l.Error("auth.renderAuthForm template parse error",
			slog.String("err", err.Error()),
			slog.String("template", tmplName))

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TemplateAuthData{
		Title: "Авторизация",
		Error: errMsg,
	}

	w.WriteHeader(statusCode)

	if err := tmpl.Execute(w, data); err != nil {
		l.Error("auth.renderAuthForm template execute error",
			slog.String("err", err.Error()),
			slog.String("template", tmplName))
	}
}
