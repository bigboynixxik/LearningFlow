package cart

import (
	"context"
	"errors"
	"html/template"
	"learningflow/internal/models"
	"learningflow/internal/transport/ssr"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
)

type SlotService interface {
	GetCart(ctx context.Context, studentID string) ([]models.Slot, error)
	AddToCart(ctx context.Context, slotID string, studentID string) error
	RemoveFromCart(ctx context.Context, slotID string, studentID string) error
	CheckoutCart(ctx context.Context, studentID string) error
}

type HandlerCart struct {
	slotSvc SlotService
}

func NewHandlerCart(slotSvc SlotService) *HandlerCart {
	return &HandlerCart{slotSvc: slotSvc}
}

type CartData struct {
	Title     string
	IsAuth    bool
	CartCount int
	Slots     []models.Slot
	TotalSum  int
	Error     string
}

func (h *HandlerCart) ShowCart(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	studentID := r.Context().Value(ssr.UserIDKey).(string)

	slots, err := h.slotSvc.GetCart(r.Context(), studentID)
	if err != nil {
		l.Error("cart.ShowCart GetCart error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalSum := 0
	for _, s := range slots {
		totalSum += s.HourlyRate
	}

	errMsg := ""
	queryErr := r.URL.Query().Get("err")
	if queryErr == "unavailable" {
		errMsg = "Извините, этот слот только что был забронирован кем-то другим."
	} else if queryErr == "empty" {
		errMsg = "Нельзя оформить пустую корзину."
	}

	data := CartData{
		Title:     "Корзина",
		IsAuth:    true,
		CartCount: len(slots),
		Slots:     slots,
		TotalSum:  totalSum,
		Error:     errMsg,
	}

	tmpl, err := template.ParseFiles("web/templates/cart.html", "web/templates/header.html")
	if err != nil {
		l.Error("cart.ShowCart template parse error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		l.Error("cart.ShowCart template execute error", slog.String("err", err.Error()))
	}
}

func (h *HandlerCart) AddToCart(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	studentID := r.Context().Value(ssr.UserIDKey).(string)
	slotID := r.FormValue("slot_id")

	err := h.slotSvc.AddToCart(r.Context(), slotID, studentID)
	if err != nil {
		l.Error("cart.AddToCart error", slog.String("err", err.Error()))
		// ПЕРЕХВАТЫВАЕМ ДОМЕННУЮ ОШИБКУ
		if errors.Is(err, models.ErrSlotUnavailable) {
			http.Redirect(w, r, "/cart?err=unavailable", http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *HandlerCart) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	studentID := r.Context().Value(ssr.UserIDKey).(string)
	slotID := r.FormValue("slot_id")

	err := h.slotSvc.RemoveFromCart(r.Context(), slotID, studentID)
	if err != nil {
		l.Error("cart.RemoveFromCart error", slog.String("err", err.Error()))
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (h *HandlerCart) Checkout(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	studentID := r.Context().Value(ssr.UserIDKey).(string)

	err := h.slotSvc.CheckoutCart(r.Context(), studentID)
	if err != nil {
		l.Error("cart.Checkout error", slog.String("err", err.Error()))
		if errors.Is(err, models.ErrCartEmpty) {
			http.Redirect(w, r, "/cart?err=empty", http.StatusSeeOther)
			return
		}
		http.Error(w, "Ошибка при оформлении", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
