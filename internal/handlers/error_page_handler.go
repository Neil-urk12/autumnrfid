package handlers

import "github.com/gofiber/fiber/v2"

func (h *AppHandler) HandleError(c *fiber.Ctx) error {
	return c.Render("partials/error", fiber.Map{})
}
