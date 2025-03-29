package main

import (
	"fmt"
	"html/template"
	"net/http"

	"admin/internal/handlers"
)

var templates *template.Template

func main() {
	mux := http.NewServeMux()
	initTemplates()

	// Routes
	mux.HandleFunc("/", handlers.IndexHandler(templates))
	mux.HandleFunc("/login", handlers.LoginPageHandler(templates))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))
	mux.HandleFunc("/logout", handlers.LogoutHandler())
	mux.HandleFunc("/students", handlers.StudentsPageHandler(templates))
	mux.HandleFunc("/dashboard", handlers.DashboardHandler(templates))

	fmt.Println("Server running on port: 8080")
	http.ListenAndServe(":8080", mux)
}

func initTemplates() {
	// Parse all templates at startup
	var err error
	templates, err = template.ParseFiles(
		"ui/html/templates/index.html",
		"ui/html/templates/login.html",
		"ui/html/partials/dashboard.html",
		"ui/html/partials/student_management.html",
	)
	if err != nil {
		// Prevent the admin from running if there are template errors
		panic(fmt.Sprintf("Failed to parse templates: %v", err))
	}
}
