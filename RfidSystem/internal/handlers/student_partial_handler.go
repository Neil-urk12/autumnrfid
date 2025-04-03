package handlers

import (
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) GetStudentPartial(ctx *fiber.Ctx) error {
	rfid := ctx.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return ctx.Render("error_page", fiber.Map{})
	}

	return ctx.Render("partials/student_info", fiber.Map{
		"Student": student,
	}, "")
}

func (h *AppHandler) HandleStudentPartial(c *fiber.Ctx) error {
	rfid := c.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return c.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	return c.Render("partials/student_info", fiber.Map{
		"Student": student,
	}, "")
}
