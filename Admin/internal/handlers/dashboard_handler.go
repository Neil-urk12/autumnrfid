package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func DashboardHandler(templates *template.Template) http.HandlerFunc {
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
			AverageScansPerDay int
		}{
			Title:              "Student Management System",
			CurrentPage:        "dashboard",
			TotalStudents:      1220,
			TotalCourses:       45,
			AverageScansPerDay: 169,
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := templates.ExecuteTemplate(w, "dashboard", data); err != nil {
				fmt.Println("Template error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
				fmt.Println("Template error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
}
