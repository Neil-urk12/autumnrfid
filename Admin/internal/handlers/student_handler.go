package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func (h *AppHandler) StudentsPageHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Url is hit!")
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		students := []map[string]string{
			{"ID": "ST001", "Name": "Jan Doe rosa", "Email": "janros@gmail.com", "Course": "Computer Science", "Status": "Active"},
			{"ID": "ST002", "Name": "Jan Cez Casupanan", "Email": "cez@casupanan.com", "Course": "Information Science", "Status": "Active"},
		}

		data := struct {
			Students []map[string]string
			Courses  []string
		}{
			Students: students,
			Courses:  []string{"Computer Science", "Business Administration", "Electrical Engineering"},
		}

		if r.Header.Get("HX-Request") == "true" {
			if err := templates.ExecuteTemplate(w, "students", data); err != nil {
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
