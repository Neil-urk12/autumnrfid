package handlers

import (
	"fmt"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	rfid := ctx.FormValue("rfid")
	if rfid == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	if student == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	fmt.Println("Student : ", student)
	// Get student's grades
	// grades, err := h.db.GetStudentGrades(student.ID)
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).SendString("Error fetching grades")
	// }

	// Pass data to the template
	return ctx.Render("home", fiber.Map{
		"Student":   student,
		"Title":     "Student Information",
		"YearLevel": getYearLevelString(student.YearLevel),
	})
}

func getYearLevelString(year int) string {
	switch year {
	case 1:
		return "First"
	case 2:
		return "Second"
	case 3:
		return "Third"
	case 4:
		return "Fourth"
	default:
		return fmt.Sprintf("%dth", year)
	}
}
func (h *AppHandler) HandleSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	return c.SendString("data: connected\n\n")
}

// func (h *AppHandler) HandleSSEEvent(c *fiber.Ctx) error {
// 	rfid := c.FormValue("rfid")
// 	if rfid == "" {
// 		return c.Status(fiber.StatusBadRequest).SendString("RFID is required")
// 	}
