package handlers

import (
	"fmt"
	"log"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	rfid := ctx.FormValue("rfid")
	if rfid == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	rfidRepo := repositories.NewRFIDRepository(h.db)
	// student, err := rfidRepo.GetStudentByRFID(rfid)
	student, err := rfidRepo.GetStudentSummaryData(rfid)

	if err != nil {
		GetBroadcaster().Broadcast("error", fmt.Sprintf(`{"message": "Database error: %v"}`, err))
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Database error: %v", err))
	}

	if student == nil {
		htmxInstruction := fmt.Sprintf(`<div hx-get="/error" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`)
		GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Student not found: %s", rfid))
	}

	log.Printf("Student found: %s\n", student.Student.StudentID)

	htmxInstruction := fmt.Sprintf(`<div hx-get="/student-partial/%s" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)

	GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
	return ctx.SendString("Processing")
}
