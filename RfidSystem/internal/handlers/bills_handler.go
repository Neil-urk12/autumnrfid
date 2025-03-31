package handlers

import (
	"regexp"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleBills(ctx *fiber.Ctx) error {
	studentId := ctx.Query("student-id")

	if studentId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student Id is required")
	}

	match, err := regexp.MatchString(`^[A-Za-z0-9]{8,12}$`, studentId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error validating student ID")
	}
	if !match {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid student ID format. Must be 8-12 alphanumeric characters")
	}

	billsRepo := repositories.NewRFIDRepository(h.db)
	billsData, err := billsRepo.GetStudentBillsByRFID(studentId)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	return ctx.Render("partials/bills", fiber.Map{
		"Title": "Student Bills",
		"Bills": billsData,
	})
}
