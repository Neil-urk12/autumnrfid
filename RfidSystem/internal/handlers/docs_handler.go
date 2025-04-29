package handlers

import "github.com/gofiber/fiber/v2"

// HandleDocs handles HTTP requests to render the API documentation page.
func (h *AppHandler) HandleDocs(c *fiber.Ctx) error {
	return c.Render("pages/docs", nil)
}
