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

type HandlerCatalog struct {
	subSvc SubjectService
	tutSvc TutorService
}

func NewHandlerCatalog(subSvc SubjectService, tSvc TutorService) *HandlerCatalog {
	return &HandlerCatalog{
		subSvc: subSvc,
		tutSvc: tSvc,
	}
}

type HomeData struct {
	Title    string
	Subjects []models.Subject
}

type TutorsData struct {
	Title        string
	Tutors       []models.Tutor
	CategoryName string
}

type TutorData struct {
	Title string
	Tutor models.Tutor
}

func (h *HandlerCatalog) HandleHome(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())

	l.Info("ssr.HandleHome, request received on HandleHome")
	subjects, err := h.subSvc.GetAll(r.Context())
	if err != nil {
		l.Error("catalog.HandeHome Get All Subjects Error",
			slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := HomeData{Subjects: subjects, Title: "Learning Flow"}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		l.Error("ssr.HandleHome template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleHome template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleHome template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleHome template execute error", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerCatalog) HandleTutors(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	l.Info("ssr.HandleTutors, request received on HandleTutors")
	tutors, err := h.tutSvc.GetAll(r.Context())
	if err != nil {
		l.Error("catalog.HandleTutors Get All tutors Error",
			slog.String("error", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := TutorsData{Tutors: tutors, Title: "Tutors"}
	tmpl, err := template.ParseFiles("web/templates/tutors.html")

	if err != nil {
		l.Error("ssr.HandleTutors template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutors template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleTutors template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutors template execute error", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerCatalog) HandleCategory(w http.ResponseWriter, r *http.Request) {
	l := logger.FromContext(r.Context())
	l.Info("ssr.HandleCategory, request received on HandleCategory")

	idStr := r.PathValue("id")
	id, idErr := strconv.ParseInt(idStr, 10, 64)
	if idErr != nil {
		l.Error("ssr.HandleCategory parse id error",
			slog.String("err", idErr.Error()))
		http.Error(w, "ssr.HandleCategory parse id error", http.StatusInternalServerError)
		return
	}

	tutors, err := h.tutSvc.GetBySubjectID(r.Context(), id)
	if err != nil {
		l.Error("catalog.HandleCategory Get tutors by subject id",
			slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	category, err := h.subSvc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		l.Error("catalog.HandleCategory Get category by id",
			slog.String("err", err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := TutorsData{Tutors: tutors, Title: "Category", CategoryName: category.Name}
	tmpl, err := template.ParseFiles("web/templates/category.html")
	if err != nil {
		l.Error("ssr.HandleCategory template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleCategory template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleCategory template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleCategory template execute error", http.StatusInternalServerError)
		return
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

	data := TutorData{Tutor: *tutor, Title: "Tutor"}
	tmpl, err := template.ParseFiles("web/templates/tutor.html")

	if err != nil {
		l.Error("ssr.HandleTutor template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutor template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		l.Error("ssr.HandleTutor template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutor template execute error", http.StatusInternalServerError)
		return
	}
}
