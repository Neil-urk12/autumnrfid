package handlers

import (
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

var gradesCache = NewLRUCache(5, time.Hour)
var semesterGradesCache = NewLRUCache(5, time.Hour)

// HandleGrades handles HTTP requests to retrieve and display student grades for the current term.
// It expects a student ID via form value "rfid" or query parameter "student-id".
// It checks the cache, fetches data from the repository if not found,
// and renders the grades partial.
func (h *AppHandler) HandleGrades(ctx *fiber.Ctx) error {
	studentId := ctx.FormValue("rfid")
	if studentId == "" {
		studentId = ctx.Query("student-id")
	}
	if studentId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student Id is required")
	}
	// Try cache first
	if cached, found := gradesCache.Get(studentId); found {
		gradesData, ok := cached.(*model.Grades)
		if ok && gradesData != nil {
			log.Printf("[CACHE HIT] Grades for %s", studentId)
			currentTerm := gradesData.CurrentTerm
			isSecondSemesterAvailable := currentTerm.Semester != "First Semester"
			preparedGrades, gwaString := h.prepareGradesAndGWA(gradesData.Grades)
			return ctx.Render("partials/grades", fiber.Map{
				"Title":                     "Student Grades",
				"Student":                   gradesData.Student,
				"Term":                      gradesData.CurrentTerm,
				"Grades":                    preparedGrades,
				"GWA":                       gwaString,
				"SelectedSemester":          currentTerm.Semester,
				"IsSecondSemesterAvailable": isSecondSemesterAvailable,
			})
		}
	}

	// Get current term to determine which semester to show by default
	currentTerm, err := h.RFIDRepository.GetCurrentTerm()
	if err != nil {
		log.Printf("Error getting current term: %v", err)
		_ = h.db.LogScanEvent(studentId, &studentId, "grade_fetch_error", fmt.Sprintf("Error getting current term: %v", err), "", "failure")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	// Get grades data
	gradesData, err := h.RFIDRepository.GetStudentGradesByRFID(studentId)
	if err != nil {
		log.Printf("Error fetching grades for student %s: %v", studentId, err)
		_ = h.db.LogScanEvent(studentId, &studentId, "grade_fetch_error", fmt.Sprintf("Error fetching grades for student %s: %v", studentId, err), "", "failure")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	if gradesData == nil {
		log.Printf("Student %s not found", studentId)
		_ = h.db.LogScanEvent(studentId, &studentId, "grade_not_found", fmt.Sprintf("Grades not found for student %s", studentId), "", "failure")
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}
	// Store in cache
	gradesCache.Set(studentId, gradesData)

	// Determine if second semester is available
	isSecondSemesterAvailable := currentTerm.Semester != "First Semester"

	// Process grades and calculate GWA
	preparedGrades, gwaString := h.prepareGradesAndGWA(gradesData.Grades)
	_ = h.db.LogScanEvent(studentId, &studentId, "grade_fetch_success", fmt.Sprintf("Fetched %d grades", len(preparedGrades)), "", "success")

	return ctx.Render("partials/grades", fiber.Map{
		"Title":                     "Student Grades",
		"Student":                   gradesData.Student,
		"Term":                      gradesData.CurrentTerm,
		"Grades":                    preparedGrades,
		"GWA":                       gwaString,
		"SelectedSemester":          currentTerm.Semester,
		"IsSecondSemesterAvailable": isSecondSemesterAvailable,
	})
}

// HandleSemesterGrades handles HTMX requests to retrieve and display student grades for a specific semester.
// It expects the student ID as a path parameter and the semester as a query parameter.
// It checks the cache, fetches data from the repository if not found,
// and renders the grades table partial.
func (h *AppHandler) HandleSemesterGrades(ctx *fiber.Ctx) error {
	studentId := ctx.Params("studentId")
	semester := ctx.Query("semester")
	if studentId == "" || semester == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student ID and semester are required")
	}
	cacheKey := studentId + ":" + semester
	if cached, found := semesterGradesCache.Get(cacheKey); found {
		gradesData, ok := cached.(*model.Grades)
		if ok && gradesData != nil {
			log.Printf("[CACHE HIT] Semester grades for %s %s", studentId, semester)
			currentTerm := gradesData.CurrentTerm
			isSecondSemesterAvailable := currentTerm.Semester != "First Semester"
			preparedGrades, gwaString := h.prepareGradesAndGWA(gradesData.Grades)
			return ctx.Render("partials/grades-table", fiber.Map{
				"Student":                   gradesData.Student,
				"Term":                      gradesData.CurrentTerm,
				"Grades":                    preparedGrades,
				"GWA":                       gwaString,
				"SelectedSemester":          semester,
				"IsSecondSemesterAvailable": isSecondSemesterAvailable,
			})
		}
	}

	// Get current term to get the academic year and check semester availability
	currentTerm, err := h.RFIDRepository.GetCurrentTerm()
	if err != nil {
		log.Printf("Error getting current term: %v", err)
		_ = h.db.LogScanEvent(studentId, &studentId, "grade_fetch_error", fmt.Sprintf("Error getting current term: %v", err), "", "failure")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	// Get grades for the requested semester
	gradesData, err := h.RFIDRepository.GetStudentGradesByRFIDAndSemester(studentId, currentTerm.AcademicYear, semester)
	if err != nil {
		log.Printf("Error fetching grades for student %s semester %s: %v", studentId, semester, err)
		_ = h.db.LogScanEvent(studentId, &studentId, "grade_fetch_error", fmt.Sprintf("Error fetching grades for student %s semester %s: %v", studentId, semester, err), "", "failure")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	// Store in cache
	semesterGradesCache.Set(cacheKey, gradesData)

	// Determine if second semester is available
	isSecondSemesterAvailable := currentTerm.Semester != "First Semester"

	// Process grades and calculate GWA
	preparedGrades, gwaString := h.prepareGradesAndGWA(gradesData.Grades)
	_ = h.db.LogScanEvent(studentId, &studentId, "grade_fetch_success", fmt.Sprintf("Fetched %d grades", len(preparedGrades)), "", "success")

	// Render only the grades table container
	return ctx.Render("partials/grades-table", fiber.Map{
		"Student":                   gradesData.Student,
		"Term":                      gradesData.CurrentTerm,
		"Grades":                    preparedGrades,
		"GWA":                       gwaString,
		"SelectedSemester":          semester,
		"IsSecondSemesterAvailable": isSecondSemesterAvailable,
	})
}

// prepareGradesAndGWA processes grades and calculates GWA
func (h *AppHandler) prepareGradesAndGWA(grades []model.GradesRecord) ([]fiber.Map, string) {
	var preparedGrades []fiber.Map
	var totalGrade float64
	var countGrades int

	// Prepare grades with safe handling of nil values
	for _, grade := range grades {
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

		// Calculate running total for GWA
		if grade.FinalGrade != nil {
			totalGrade += *grade.FinalGrade
			countGrades++
		}
	}

	// Calculate GWA
	var gwaString string
	if countGrades > 0 {
		gwa := totalGrade / float64(countGrades)
		gwaString = fmt.Sprintf("%.2f", gwa)
	} else {
		gwaString = "N/A"
	}

	return preparedGrades, gwaString
}

// formatGrade safely handles nil grade values and formats non-nil grades
// Returns:
//   - "-" when the grade is nil (grade not yet available/recorded)
//   - formatted string with 2 decimal places for actual grade values
func formatGrade(grade *float64) string {
	if grade == nil || *grade == 0 {
		// Return a dash placeholder for nil grades
		// This indicates that the grade is not yet available or has not been recorded
		return "-"
	}
	return fmt.Sprintf("%.2f", *grade)
}
