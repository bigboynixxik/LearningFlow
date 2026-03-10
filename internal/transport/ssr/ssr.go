package ssr

import (
	"html/template"
	"learningflow/internal/models"
	"log/slog"
	"net/http"
)

type HomeData struct {
	Title    string
	Subjects []models.Subject
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	subjects := []models.Subject{
		{ID: 1, Name: "Математика"},
		{ID: 2, Name: "Программирование на Go"},
		{ID: 3, Name: "Английский язык"},
		{ID: 4, Name: "Физика"},
	}

	data := HomeData{Subjects: subjects, Title: "Learning Flow"}

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		slog.Error("ssr.HandleHome template error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleHome template error", http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, data); err != nil {
		slog.Error("ssr.HandleHome template error",
			slog.String("err", err.Error()))
		http.Error(w, "ssr.HandleHome template error", http.StatusInternalServerError)
	}
}
