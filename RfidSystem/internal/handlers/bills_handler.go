package handlers

import "github.com/gofiber/fiber/v2"

func (h *AppHandler) HandleBills(ctx *fiber.Ctx) error {
	studentID := ctx.Query("student-id")

	if studentID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student Id is required")
	}

	return nil
}
