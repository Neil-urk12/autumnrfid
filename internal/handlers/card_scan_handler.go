package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	rfid := ctx.FormValue("rfid")
	if rfid == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	fmt.Println("RFID:", rfid)

	student, err := h.db.GetStudentByRFID(rfid)
	// student, err := h.db.GetStudentByRFID("ACLC-2023-001")
	if err != nil {
		fmt.Println(err)
		// fmt.Println("Student", student)
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
