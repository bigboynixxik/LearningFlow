package ssr

import (
	"html/template"
	"learningflow/internal/models"
	"learningflow/pkg/logger"
	"log/slog"
	"net/http"
	"strconv"
)

// Временные данные для Лабы №1 (имитация базы данных)
var mockSubjects = []models.Subject{
	{ID: 1, Name: "Математика"},
	{ID: 2, Name: "Информатика"},
	{ID: 3, Name: "Английский язык"},
	{ID: 4, Name: "Физика"},
}

var mockTutors = []models.Tutor{
	{
		UserID:      "1",
		Name:        "Иван Иванов",
		HourlyRate:  1500,
		Description: "Информатика",
		SubjectIDs:  []int64{2},
	},
	{
		UserID:      "2",
		Name:        "Петр Петров",
		HourlyRate:  2000,
		Description: "Математика любой сложности. Подготовка к ЕГЭ.",
		SubjectIDs:  []int64{1, 4},
	},
}

type HomeData struct {
	Title    string
	Subjects []models.Subject
}

type TutorsData struct {
	Title  string
	Tutors []models.Tutor
}

type TutorData struct {
	Title string
	Tutor models.Tutor
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	log.Info("ssr.HandleHome, request received on HandleHome")
	subjects := mockSubjects

	data := HomeData{Subjects: subjects, Title: "Learning Flow"}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Error("ssr.HandleHome template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleHome template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Error("ssr.HandleHome template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleHome template execute error", http.StatusInternalServerError)
		return
	}
}

func HandleTutors(w http.ResponseWriter, r *http.Request) {
	logs := logger.FromContext(r.Context())
	logs.Info("ssr.HandleTutors, request received on HandleTutors")
	tutors := mockTutors
	data := TutorsData{Tutors: tutors, Title: "Tutors"}
	tmpl, err := template.ParseFiles("web/templates/tutors.html")

	if err != nil {
		logs.Error("ssr.HandleTutors template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutors template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		logs.Error("ssr.HandleTutors template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutors template execute error", http.StatusInternalServerError)
		return
	}
}

func HandleCategory(w http.ResponseWriter, r *http.Request) {
	logs := logger.FromContext(r.Context())
	logs.Info("ssr.HandleCategory, request received on HandleCategory")

	idStr := r.PathValue("id")
	id, idErr := strconv.ParseInt(idStr, 10, 64)
	if idErr != nil {
		logs.Error("ssr.HandleCategory parse id error",
			slog.String("err", idErr.Error()))
		http.Error(w, "ssr.HandleCategory parse id error", http.StatusInternalServerError)
		return
	}
	var filteredTutors []models.Tutor
	for _, tutor := range mockTutors {
		for _, subjectID := range tutor.SubjectIDs {
			if id == subjectID {
				filteredTutors = append(filteredTutors, tutor)
				break
			}
		}
	}
	data := TutorsData{Tutors: filteredTutors, Title: "Category"}
	tmpl, err := template.ParseFiles("web/templates/category.html")
	if err != nil {
		logs.Error("ssr.HandleCategory template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleCategory template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		logs.Error("ssr.HandleCategory template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleCategory template execute error", http.StatusInternalServerError)
		return
	}
}

func HandleTutor(w http.ResponseWriter, r *http.Request) {
	logs := logger.FromContext(r.Context())
	logs.Info("ssr.HandleTutor, request received on HandleTutor")

	id := r.PathValue("id")
	var foundTutor *models.Tutor
	for _, tutor := range mockTutors {
		if id == tutor.UserID {
			t := tutor
			foundTutor = &t
			break
		}
	}
	if foundTutor == nil {
		http.NotFound(w, r)
		return
	}

	data := TutorData{Tutor: *foundTutor, Title: "Tutor"}
	tmpl, err := template.ParseFiles("web/templates/tutor.html")

	if err != nil {
		logs.Error("ssr.HandleTutor template parse error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutor template parse error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		logs.Error("ssr.HandleTutor template execute error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleTutor template execute error", http.StatusInternalServerError)
		return
	}
}
