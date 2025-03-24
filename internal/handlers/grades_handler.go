package handlers

import (
	"fmt"
	"regexp"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleGrades(ctx *fiber.Ctx) error {
	studentId := ctx.Query("student-id")

	if studentId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student Id is required")
	}

	// Validate student ID format (assuming it should be alphanumeric and 8-12 characters)
	// You can adjust the regex as needed
	match, err := regexp.MatchString(`^[A-Za-z0-9]{8,12}$`, studentId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error validating student ID")
	}
	if !match {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid student ID format. Must be 8-12 alphanumeric characters")
	}

	gradesRepo := repositories.NewRFIDRepository(h.db)
	gradesData, err := gradesRepo.GetStudentGradesByRFID(studentId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	if gradesData == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	// err := ctx.Render("pages/grades", fiber.Map{
	// 	"Title": "Student Grades",
	// })
	return ctx.Render("grades", fiber.Map{
		"Title":  "Student Grades",
		"Grades": gradesData,
	})

	// if err != nil {
	// 	fmt.Printf("Template rendering error: %v\n", err)
	// 	return err
	// }
}

func (h *AppHandler) HandleTestGrades(ctx *fiber.Ctx) error {
	err := ctx.Render("partials/grades", fiber.Map{
		"Title": "Student Grades",
	})
	if err != nil {
		fmt.Printf("Template rendering error: %v\n", err)
		return err
	}
	return nil
}
