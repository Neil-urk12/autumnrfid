package handlers

import (
	"fmt"
	"log"
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
	// match, err := regexp.MatchString(`^[A-Za-z0-9]{8,12}$`, studentId)
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).SendString("Error validating student ID")
	// }
	// if !match {
	// 	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid student ID format. Must be 8-12 alphanumeric characters")
	// }
	log.Printf("Received request for student ID: %s", studentId)

	gradesRepo := repositories.NewRFIDRepository(h.db)
	gradesData, err := gradesRepo.GetStudentGradesByRFID(studentId)
	if err != nil {
		log.Printf("Error fetching grades for student %s: %v", studentId, err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	if gradesData == nil {
		log.Printf("Student %s not found", studentId)
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	var preparedGrades []fiber.Map

	for _, grade := range gradesData.Grades {
		preparedGrades = append(preparedGrades, fiber.Map{
			"SubjectCode":    grade.SubjectCode,
			"SubjectName":    grade.SubjectName,
			"PrelimGrade":    formatGrade(grade.PrelimGrade),
			"MidtermGrade":   formatGrade(grade.MidtermGrade),
			"PrefinalGrade":  formatGrade(grade.PrefinalGrade),
			"FinalTermGrade": formatGrade(grade.FinalTermGrade),
			"FinalGrade":     formatGrade(grade.FinalGrade),
		})
	}

	var totalGrade float64
	var countGrades int
	for _, grade := range gradesData.Grades {
		if grade.FinalGrade != nil {
			log.Printf("Grade for subject %s: %.2f", grade.SubjectCode, *grade.FinalGrade)
			totalGrade += *grade.FinalGrade
			countGrades++
		}
	}

	var gwa float64
	if countGrades > 0 {
		gwa = totalGrade / float64(countGrades)
	}

	log.Printf("Student %s has %d grades with a total of %.2f and a GWA of %.2f", studentId, countGrades, totalGrade, gwa)

	return ctx.Render("partials/grades", fiber.Map{
		"Title":  "Student Grades",
		"Term":   gradesData.CurrentTerm,
		"Grades": gradesData.Grades,
		"GWA":    fmt.Sprintf("%.2f", gwa),
	})
}

func formatGrade(grade *float64) string {
	if grade == nil {
		return "-" // Placeholder for nil values
	}
	return fmt.Sprintf("%.2f", *grade) // Dereference and format
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
