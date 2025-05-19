package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var users = map[string]string{
	"admin@gmail.com": "admin123",
	"user@gmail.com":  "user123",
}

func (h *AppHandler) LoginPageHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == fiber.MethodGet {
			return c.Render("pages/login", fiber.Map{
				"RedirectURL": c.Query("redirect"),
			})
		}

		if c.Method() == fiber.MethodPost {
			email := c.FormValue("email")
			password := c.FormValue("password")
			redirectURL := c.FormValue("redirect_url", "/logs")
			if authenticated, err := authenticateUser(email, password); authenticated {
				sessionToken := CreateSession(email)

				if c.Get("HX-Request") == "true" {
					c.Set("HX-Redirect", redirectURL)
					c.Set("HX-Trigger", `{"loginSuccessClient": {"token": "`+sessionToken+`"}}`)
					return c.SendStatus(fiber.StatusOK)
				} else {
					return c.Redirect(redirectURL, fiber.StatusSeeOther)
				}
			} else {
				errorMsg := "Invalid email or password"
				if err != nil {
					errorMsg = err.Error()
				}
				return c.Render("pages/login", fiber.Map{
					"Error":       errorMsg,
					"RedirectURL": redirectURL,
				})
			}
		}
		return errors.New("method not allowed for /login")
	}
}

func (h *AppHandler) LogoutHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() != fiber.MethodPost {
			return c.Status(fiber.StatusMethodNotAllowed).SendString("Method Not Allowed")
		}

		authHeader := c.Get("Authorization")
		token := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token != "" {
			DeleteSession(token)
		}

		if c.Get("HX-Request") == "true" {
			c.Set("HX-Trigger", `{"logoutSuccessClient": {}}`)
			c.Set("HX-Redirect", "/login")
			return c.SendStatus(fiber.StatusOK)
		} else {
			return c.Redirect("/login", fiber.StatusSeeOther)
		}
	}
}

func authenticateUser(email, password string) (bool, error) {
	if email == "" || password == "" {
		return false, errors.New("email and password are required")
	}
	storedPassword, exists := users[email]
	if !exists {
		return false, errors.New("invalid email or password")
	}
	if storedPassword != password {
		return false, errors.New("invalid email or password")
	}
	return true, nil
}

// AuthRequired is a middleware that checks if the user is authenticated.
// If not, it handles the response based on request type (HTMX or regular).
// func AuthRequired() fiber.Handler { // THIS FUNCTION WILL BE REMOVED
// 	return func(c *fiber.Ctx) error {
// 		if !IsAuthenticatedFiber(c) {
// 			originalURL := c.OriginalURL()
// 			if originalURL == "/login" || originalURL == "/logout" || originalURL == "/" || originalURL == "/favicon.ico" {
// 				originalURL = "/logs"
// 			}
// 			if c.Get("HX-Request") == "true" {
// 				c.Set("HX-Redirect", "/login?redirect="+originalURL)
// 				return c.Status(fiber.StatusUnauthorized).SendString("HTMX redirecting to login")
// 			} else {
// 				return c.Redirect("/login?redirect="+originalURL, fiber.StatusSeeOther)
// 			}
// 		}
// 		return c.Next()
// 	}
// }

// IsAuthenticatedFiber, GetSessionUserEmailFiber, CreateSession, DeleteSession
// might need internal adjustments if they relied solely on cookies previously.
// For now, their direct calls from the modified handlers assume CreateSession returns a token
// and DeleteSession invalidates a token. IsAuthenticatedFiber is no longer called by AuthRequired.
// If other handlers need to check authentication, they will need to use a revised IsAuthenticatedFiber
// that checks the Authorization header.
