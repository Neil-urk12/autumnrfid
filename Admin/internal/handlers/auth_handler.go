package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var users = map[string]string{
	"admin@example.com": "admin123",
	"user@example.com":  "user123",
}

func LoginPageHandler(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// If already logged in, redirect to dashboard
		if IsAuthenticated(r) {
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
				sessionToken := CreateSession(email)

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
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear session cookie
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err == nil {
			sessionToken := cookie.Value
			DeleteSession(sessionToken)
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
