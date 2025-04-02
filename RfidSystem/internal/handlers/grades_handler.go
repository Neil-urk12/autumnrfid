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

	// Prepare grades with safe handling of nil values
	for _, grade := range gradesData.Grades {
		// preparedGrades = append(preparedGrades, fiber.Map{
		// 	"SubjectCode":    grade.SubjectCode,
		// 	"SubjectName":    grade.SubjectName,
		// 	"PrelimGrade":    formatGrade(grade.PrelimGrade),
		// 	"MidtermGrade":   formatGrade(grade.MidtermGrade),
		// 	"PrefinalGrade":  formatGrade(grade.PrefinalGrade),
		// 	"FinalTermGrade": formatGrade(grade.FinalTermGrade),
		// 	"FinalGrade":     formatGrade(grade.FinalGrade),
		// })
		gradeMap := fiber.Map{
			"SubjectCode":    grade.SubjectCode,
			"SubjectName":    grade.SubjectName,
			"PrelimGrade":    formatGrade(grade.PrelimGrade),
			"MidtermGrade":   formatGrade(grade.MidtermGrade),
			"PrefinalGrade":  formatGrade(grade.PrefinalGrade),
			"FinalTermGrade": formatGrade(grade.FinalTermGrade),
			"FinalGrade":     formatGrade(grade.FinalGrade),
		}
		preparedGrades = append(preparedGrades, gradeMap)
	}

	// Calculate GWA with nil safety checks
	var totalGrade float64
	var countGrades int
	for _, grade := range gradesData.Grades {
		if grade.FinalGrade != nil {
			totalGrade += *grade.FinalGrade
			countGrades++
			log.Printf("Grade for subject %s: %.2f", grade.SubjectCode, *grade.FinalGrade)
		}
	}

	var gwa float64
	var gwaString string
	if countGrades > 0 {
		gwa = totalGrade / float64(countGrades)
		gwaString = fmt.Sprintf("%.2f", gwa)
	} else {
		// Use a placeholder dash when no grades are available for GWA calculation
		gwaString = "N/A"
	}

	log.Printf("Student %s has %d grades with a total of %.2f and a GWA of %.2f", studentId, countGrades, totalGrade, gwa)

	return ctx.Render("partials/grades", fiber.Map{
		"Title": "Student Grades",
		"Term":  gradesData.CurrentTerm,
		// "Grades": gradesData.Grades,
		"Grades": preparedGrades,
		"GWA":    gwaString,
	})
}

// formatGrade safely handles nil grade values and formats non-nil grades
// Returns:
//   - "-" when the grade is nil (grade not yet available/recorded)
//   - formatted string with 2 decimal places for actual grade values
func formatGrade(grade *float64) string {
	if grade == nil {
		// Return a dash placeholder for nil grades
		// This indicates that the grade is not yet available or has not been recorded
		return "-"
	}
	return fmt.Sprintf("%.2f", *grade)
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
