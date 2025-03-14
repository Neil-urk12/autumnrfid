package handlers

import "github.com/gofiber/fiber/v2"

func (h *AppHandler) HandleGrades(ctx *fiber.Ctx) error {
	return ctx.Render("grades", fiber.Map{})
}
