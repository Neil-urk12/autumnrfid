package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type session struct {
	email  string
	expiry time.Time
}

var sessions = map[string]session{}

// CreateSession creates a new session for the given email and returns the session token.
func CreateSession(email string) string {
	sessionToken := generateSessionToken()
	expiry := time.Now().Add(24 * time.Hour)

	sessions[sessionToken] = session{
		email:  email,
		expiry: expiry,
	}

	return sessionToken
}

// DeleteSession removes a session from the sessions map.
func DeleteSession(sessionToken string) {
	delete(sessions, sessionToken)
}

// IsAuthenticated checks if a user is authenticated based on the session cookie in a standard http.Request.
// This function is kept for backward compatibility.
func IsAuthenticated(r *http.Request) bool {
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

// IsAuthenticatedFiber checks if a user is authenticated based on the session cookie using a Fiber context.
func IsAuthenticatedFiber(c *fiber.Ctx) bool {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return false
	}

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

// GetSessionUserEmailFiber retrieves the email address of the logged-in user from their session.
// Returns the email address and true if a valid session exists, empty string and false otherwise.
func GetSessionUserEmailFiber(c *fiber.Ctx) (string, bool) {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return "", false
	}

	userSession, exists := sessions[sessionToken]
	if !exists {
		return "", false
	}

	// Check if session has expired
	if userSession.expiry.Before(time.Now()) {
		delete(sessions, sessionToken)
		return "", false
	}

	return userSession.email, true
}

func generateSessionToken() string {
	// Simple implementation for demonstration
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
