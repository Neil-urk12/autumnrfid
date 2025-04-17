package handlers

import (
	"admin/repositories"
	"fmt"
	"html/template"
	"net/http"
)

type AppHandler struct {
	db *repositories.DatabaseClient
}

func NewHandler(db *repositories.DatabaseClient) *AppHandler {
	return &AppHandler{db: db}
}

func (h *AppHandler) IndexHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		data := struct {
			Title              string
			CurrentPage        string
			TotalStudents      int
			TotalCourses       int
			TotalBills         int
			AverageScansPerDay int
		}{
			Title:              "Student Management System",
			CurrentPage:        "dashboard",
			TotalStudents:      1220,
			TotalCourses:       45,
			AverageScansPerDay: 169,
		}
		if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
			fmt.Println("Template error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		fmt.Println("User authenticated, serving index page")
	}
}
