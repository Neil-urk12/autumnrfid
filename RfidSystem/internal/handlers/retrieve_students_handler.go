package handlers

import (
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) RetrieveStudentsHandler(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)

	rfidRepo := repositories.NewRFIDRepository(h.db)
	students, err := rfidRepo.GetAllStudents(page)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return ctx.JSON(students)
}
