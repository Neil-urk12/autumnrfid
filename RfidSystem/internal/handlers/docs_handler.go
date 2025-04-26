package handlers

import "github.com/gofiber/fiber/v2"

func (h *AppHandler) HandleDocs(c *fiber.Ctx) error {
	return c.Render("pages/docs", nil)
}
