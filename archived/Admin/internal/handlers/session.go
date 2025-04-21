package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type session struct {
	email  string
	expiry time.Time
}

var sessions = map[string]session{}

func CreateSession(email string) string {
	sessionToken := generateSessionToken()
	expiry := time.Now().Add(24 * time.Hour)

	sessions[sessionToken] = session{
		email:  email,
		expiry: expiry,
	}

	return sessionToken
}

func DeleteSession(sessionToken string) {
	delete(sessions, sessionToken)
}

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

func generateSessionToken() string {
	// Simple implementation for demonstration
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
