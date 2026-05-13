package profile

import (
	"context"
	"errors"
	"html/template"
	"learningflow/internal/models"
	"learningflow/internal/transport/ssr"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
	"time"
)

type SlotService interface {
	CreateAvailability(ctx context.Context, tutorID string, startTime time.Time) error
	GetCartCount(ctx context.Context, studentID string) (int, error)
}

type HandlerProfile struct {
	slotSvc SlotService
}

func NewHandlerProfile(slotSvc SlotService) *HandlerProfile {
	return &HandlerProfile{slotSvc: slotSvc}
}

type ProfileData struct {
	Title     string
	IsAuth    bool
	CartCount int
	Error     string
	Success   string
}

func (h *HandlerProfile) ShowProfile(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	userID := r.Context().Value(ssr.UserIDKey).(string)

	count, _ := h.slotSvc.GetCartCount(r.Context(), userID)

	errMsg := r.URL.Query().Get("err")
	successMsg := r.URL.Query().Get("success")

	if errMsg == "overlap" {
		errMsg = "У вас уже есть занятие на это время."
	} else if errMsg == "past" {
		errMsg = "Нельзя создать занятие в прошлом."
	} else if errMsg == "invalid" {
		errMsg = "Неверный формат времени."
	} else if successMsg == "ok" {
		successMsg = "Слот успешно добавлен в расписание!"
	}

	data := ProfileData{
		Title:     "Личный кабинет",
		IsAuth:    true,
		CartCount: count,
		Error:     errMsg,
		Success:   successMsg,
	}

	tmpl, err := template.ParseFiles("web/templates/profile.html", "web/templates/header.html")
	if err != nil {
		l.Error("profile.ShowProfile template parse error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		l.Error("profile.ShowProfile template execute error", slog.String("err", err.Error()))
	}
}

func (h *HandlerProfile) AddSlot(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	userID := r.Context().Value(ssr.UserIDKey).(string)

	startTimeStr := r.FormValue("start_time")

	startTime, err := time.Parse("2006-01-02T15:04", startTimeStr)
	if err != nil {
		l.Error("profile.AddSlot parse time error", slog.String("err", err.Error()))
		http.Redirect(w, r, "/profile?err=invalid", http.StatusSeeOther)
		return
	}

	err = h.slotSvc.CreateAvailability(r.Context(), userID, startTime)
	if err != nil {
		l.Error("profile.AddSlot create error", slog.String("err", err.Error()))

		if errors.Is(err, models.ErrSlotPast) {
			http.Redirect(w, r, "/profile?err=past", http.StatusSeeOther)
			return
		}
		if errors.Is(err, models.ErrSlotOverlap) {
			http.Redirect(w, r, "/profile?err=overlap", http.StatusSeeOther)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile?success=ok", http.StatusSeeOther)
}
