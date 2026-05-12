package catalog

import (
	"context"
	"errors"
	"html/template"
	"learningflow/internal/models"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
	"strconv"
)

type TutorService interface {
	GetAll(ctx context.Context) ([]models.Tutor, error)
	GetByID(ctx context.Context, id string) (*models.Tutor, error)
	GetBySubjectID(ctx context.Context, subjectID int64) ([]models.Tutor, error)
}

type SubjectService interface {
	GetAll(ctx context.Context) ([]models.Subject, error)
	GetByID(ctx context.Context, id int64) (*models.Subject, error)
}

type SlotService interface {
	GetFreeSlots(ctx context.Context, tutorID string) ([]models.Slot, error)
	GetCartCount(ctx context.Context, studentID string) (int, error)
}

type SessionValidator interface {
	ValidateSession(ctx context.Context, token string) (string, error)
}

type HandlerCatalog struct {
	subSvc  SubjectService
	tutSvc  TutorService
	slotSvc SlotService
	authSvc SessionValidator
}

func NewHandlerCatalog(subSvc SubjectService, tSvc TutorService, slotSvc SlotService, authSvc SessionValidator) *HandlerCatalog {
	return &HandlerCatalog{
		subSvc:  subSvc,
		tutSvc:  tSvc,
		slotSvc: slotSvc,
		authSvc: authSvc,
	}
}

type HomeData struct {
	Title     string
	Subjects  []models.Subject
	IsAuth    bool
	CartCount int
}

type TutorsData struct {
	Title        string
	Tutors       []models.Tutor
	CategoryName string
	IsAuth       bool
	CartCount    int
}

type TutorData struct {
	Title     string
	Tutor     models.Tutor
	FreeSlots []models.Slot
	IsAuth    bool
	CartCount int
}

func (h *HandlerCatalog) getAuthAndCartCount(ctx context.Context, r *http.Request) (bool, int) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, 0
	}

	userID, err := h.authSvc.ValidateSession(ctx, cookie.Value)
	if err != nil {
		return false, 0 // ИСПРАВЛЕНО: добавили 0
	}

	count, _ := h.slotSvc.GetCartCount(ctx, userID)
	return true, count
}

func (h *HandlerCatalog) HandleHome(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	l.Info("ssr.HandleHome, request received on HandleHome")

	subjects, err := h.subSvc.GetAll(r.Context())
	if err != nil {
		l.Error("catalog.HandleHome Get All Subjects Error", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAuth, cartCount := h.getAuthAndCartCount(r.Context(), r)

	data := HomeData{
		Subjects:  subjects,
		Title:     "Learning Flow",
		IsAuth:    isAuth,
		CartCount: cartCount,
	}

	tmpl, err := template.ParseFiles("web/templates/index.html", "web/templates/header.html")
	if err != nil {
		l.Error("ssr.HandleHome template parse error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleHome template execute error", slog.String("err", err.Error()))
	}
}

func (h *HandlerCatalog) HandleTutors(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	l.Info("ssr.HandleTutors, request received on HandleTutors")

	tutors, err := h.tutSvc.GetAll(r.Context())
	if err != nil {
		l.Error("catalog.HandleTutors Get All tutors Error", slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAuth, cartCount := h.getAuthAndCartCount(r.Context(), r)

	data := TutorsData{
		Tutors:    tutors,
		Title:     "Tutors",
		IsAuth:    isAuth,
		CartCount: cartCount,
	}

	tmpl, err := template.ParseFiles("web/templates/tutors.html", "web/templates/header.html")
	if err != nil {
		l.Error("ssr.HandleTutors template parse error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleTutors template execute error", slog.String("err", err.Error()))
	}
}

func (h *HandlerCatalog) HandleCategory(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	l.Info("ssr.HandleCategory, request received on HandleCategory")

	idStr := r.PathValue("id")
	id, idErr := strconv.ParseInt(idStr, 10, 64)
	if idErr != nil {
		l.Error("ssr.HandleCategory parse id error", slog.String("err", idErr.Error()))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tutors, err := h.tutSvc.GetBySubjectID(r.Context(), id)
	if err != nil {
		l.Error("catalog.HandleCategory Get tutors by subject id", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	category, err := h.subSvc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		l.Error("catalog.HandleCategory Get category by id", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	isAuth, cartCount := h.getAuthAndCartCount(r.Context(), r)

	data := TutorsData{
		Tutors:       tutors,
		Title:        "Category",
		CategoryName: category.Name,
		IsAuth:       isAuth,
		CartCount:    cartCount,
	}

	tmpl, err := template.ParseFiles("web/templates/category.html", "web/templates/header.html")
	if err != nil {
		l.Error("ssr.HandleCategory template parse error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleCategory template execute error", slog.String("err", err.Error()))
	}
}

func (h *HandlerCatalog) HandleTutor(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	l.Info("ssr.HandleTutor, request received on HandleTutor")

	id := r.PathValue("id")
	tutor, err := h.tutSvc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		l.Error("catalog.HandleTutor Get tutor by id", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	freeSlots, err := h.slotSvc.GetFreeSlots(r.Context(), id)
	if err != nil {
		l.Error("catalog.HandleTutor GetFreeSlots error", slog.String("err", err.Error()))
	}

	isAuth, cartCount := h.getAuthAndCartCount(r.Context(), r)

	data := TutorData{
		Tutor:     *tutor,
		Title:     tutor.Name,
		FreeSlots: freeSlots,
		IsAuth:    isAuth,
		CartCount: cartCount,
	}

	tmpl, err := template.ParseFiles("web/templates/tutor.html", "web/templates/header.html")
	if err != nil {
		l.Error("ssr.HandleTutor template parse error", slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleTutor template execute error", slog.String("err", err.Error()))
	}
}
