package handlers

import "github.com/gofiber/fiber/v2"

func (h *AppHandler) HandleError(c *fiber.Ctx) error {
	return c.Render("pages/error_page", fiber.Map{
		"Title": "Error",
	})
}
