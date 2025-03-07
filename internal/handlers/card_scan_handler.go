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

	fmt.Println("RFID:", rfid)

	// Create a new RFIDRepository instance
	rfidRepo := repositories.NewRFIDRepository(h.db)

	student, err := rfidRepo.GetStudentByRFID(rfid)
	if err != nil {
		fmt.Println(err)
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
		"Student": student,
		// "Grades":  grades,
		// "Bills":   bills,
	})
}
