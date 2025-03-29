package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginPageHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))
	mux.HandleFunc("/logout", logoutHandler)

	fmt.Println("Server running on port: 8080")
	http.ListenAndServe(":8080", mux)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

var users = map[string]string{
	"admin@example.com": "admin123",
	"user@example.com":  "user123",
}
var sessions = map[string]session{}

type session struct {
	email  string
	expiry time.Time
}

func createSession(email string) string {
	sessionToken := generateSessionToken()
	expiry := time.Now().Add(24 * time.Hour)

	sessions[sessionToken] = session{
		email:  email,
		expiry: expiry,
	}

	return sessionToken
}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	sessionToken := cookie.Value
	userSession, exists := sessions[sessionToken]
	if !exists {
		return false
	}

	// Check if session has expired
	if userSession.expiry.Before(time.Now()) {
		delete(sessions, sessionToken)
		return false
	}

	return true
}

func generateSessionToken() string {
	// Simple implementation for demonstration
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := struct {
		Title              string
		TotalStudents      int
		TotalCourses       int
		TotalBills         int
		AverageScansPerDay int
	}{
		Title:              "Student Management System",
		TotalStudents:      1220,
		TotalCourses:       45,
		AverageScansPerDay: 169,
	}
	tmpl := template.Must(template.ParseFiles("ui/html/templates/index.html"))
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Println(err)
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	// If already logged in, redirect to dashboard
	if isAuthenticated(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("ui/html/templates/login.html"))
		if err := tmpl.Execute(w, nil); err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

	case http.MethodPost:
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get email and password from form
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Check if credentials are valid
		if authenticated, err := authenticateUser(email, password); authenticated {
			// Create session
			sessionToken := createSession(email)

			// Set cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    sessionToken,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(24 * time.Hour),
			})

			// Check if this is an HTMX request
			if r.Header.Get("HX-Request") == "true" {
				// Use HX-Redirect for HTMX requests
				w.Header().Set("HX-Redirect", "/")
				w.WriteHeader(http.StatusOK)
				return
			} else {
				// Regular redirect for non-HTMX requests
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		} else {
			// Authentication failed
			errorMsg := "Invalid email or password"
			if err != nil {
				errorMsg = err.Error()
			}

			// Prepare error data
			data := struct {
				Error string
			}{
				Error: errorMsg,
			}

			// Render the login template with error message
			tmpl := template.Must(template.ParseFiles("ui/html/templates/login.html"))
			if err := tmpl.Execute(w, data); err != nil {
				fmt.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
}

func authenticateUser(email, password string) (bool, error) {
	if email == "" || password == "" {
		return false, errors.New("email and password are required")
	}

	// Check if user exists
	storedPassword, exists := users[email]
	if !exists {
		return false, errors.New("invalid email or password")
	}

	// Check if password matches
	if storedPassword != password {
		return false, errors.New("invalid email or password")
	}

	return true, nil
}

type LoginDetails struct {
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear session cookie
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil {
		sessionToken := cookie.Value
		delete(sessions, sessionToken)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
	})

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		// Use HX-Redirect for HTMX requests to do a client-side redirect
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	} else {
		// Regular redirect for non-HTMX requests
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
