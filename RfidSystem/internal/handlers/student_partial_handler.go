package handlers

import (
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

var studentInfoCache = NewLRUCache(5, time.Hour)

// HandleStudentInfo handles HTTP requests to render the student information partial.
// It supports receiving the student ID via POST body or GET path parameter.
// It checks the cache, fetches student summary data from the repository if necessary,
// stores the data in the cache, and renders the student info partial.
func (h *AppHandler) HandleStudentInfo(ctx *fiber.Ctx) error {
	// Support POST body or GET path parameter
	var req struct {
		RFID string `json:"rfid" form:"rfid"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		if err2 := h.db.LogScanEvent(req.RFID, nil, "student_info_error", fmt.Sprintf("Error parsing request body: %v", err), "", "failure"); err2 != nil {
			log.Printf("Failed to log student_info_error: %v", err2)
		}
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	studentId := req.RFID
	if studentId == "" {
		studentId = ctx.Params("rfid")
	}
	if studentId == "" {
		log.Printf("Student ID is required")
		_ = h.db.LogScanEvent("", nil, "student_info_error", "Student ID is required", "", "failure")
		return ctx.Status(fiber.StatusBadRequest).SendString("Student ID is required")
	}

	// Try cache first
	if cached, found := studentInfoCache.Get(studentId); found {
		studentInfo, ok := cached.(*model.StudentInfoViewModel)
		if ok && studentInfo != nil && studentInfo.Student != nil {
			// Update last access timestamp directly
			now := time.Now()
			updateQuery := `
			UPDATE Students
			SET last_access_timestamp = ?
			WHERE student_ID = ?`

			if _, err := h.db.DB.Exec(updateQuery, now, studentId); err != nil {
				log.Printf("Error updating access timestamps for student %s: %v", studentId, err)
			}
			// Update last access timestamp even for cached data
			var formattedSchedules []model.PaymentScheduleViewModel
			for _, schedule := range studentInfo.PaymentSchedules {
				if schedule.TermDescription == "Initial Payment" {
					continue
				}
				formattedSchedules = append(formattedSchedules, model.PaymentScheduleViewModel{
					ID:                      schedule.ID,
					AssessmentNumber:        schedule.AssessmentNumber,
					TermDescription:         schedule.TermDescription,
					DueDate:                 schedule.DueDate,
					ExpectedAmount:          schedule.ExpectedAmount,
					ExpectedAmountFormatted: formatAmount(schedule.ExpectedAmount),
					SortOrder:               schedule.SortOrder,
				})
			}
			_ = h.db.LogScanEvent(studentId, &studentInfo.Student.StudentID, "info_displayed", fmt.Sprintf("Displayed cached info for student ID %s", studentId), "", "success")
			return ctx.Render("partials/student_info", fiber.Map{
				"Student":          studentInfo.Student,
				"YearLevel":        studentInfo.YearLevel,
				"GradesSummary":    studentInfo.GradesSummary,
				"Assessment":       formatAssessmentForView(studentInfo.Assessment),
				"PaymentSchedules": formattedSchedules,
			})
		}
	}

	studentInfo, err := h.RFIDRepository.GetStudentSummaryData(studentId)
	if err != nil {
		log.Printf("Error getting student info: %v", err)
		_ = h.db.LogScanEvent(studentId, nil, "db_error", fmt.Sprintf("Error getting student info: %v", err), "", "failure")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve student information")
	}

	if studentInfo == nil || studentInfo.Student == nil {
		_ = h.db.LogScanEvent(studentId, nil, "student_not_found", fmt.Sprintf("Student not found: %s", studentId), "", "failure")
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}
	log.Printf("Student info retrieved successfully")
	// Store in cache
	studentInfoCache.Set(studentId, studentInfo)
	_ = h.db.LogScanEvent(studentId, &studentInfo.Student.StudentID, "info_displayed", fmt.Sprintf("Displayed info for student ID %s", studentId), "", "success")

	// Log grades summary for debugging
	if studentInfo.GradesSummary != nil {
		log.Printf("Found %d year summaries", len(studentInfo.GradesSummary))
		for _, year := range studentInfo.GradesSummary {
			log.Printf("Year: %s, First: %v, Second: %v",
				year.YearName,
				year.FirstSem,
				year.SecondSem)
		}
	} else {
		log.Printf("No grades summary found")
	}

	log.Printf("Assessment data: %+v", studentInfo.Assessment)

	// Format the assessment data
	var formattedAssessment model.AssessmentViewModel
	if studentInfo.Assessment != nil {
		formattedAssessment = formatAssessmentForView(studentInfo.Assessment)
	} else {
		log.Printf("No assessment data to format, setting formattedAssessment to zero value")
		formattedAssessment = model.AssessmentViewModel{}
	}
	log.Printf("Formatted assessment: %+v", formattedAssessment)

	// Format payment schedules
	var formattedSchedules []model.PaymentScheduleViewModel
	for _, schedule := range studentInfo.PaymentSchedules {
		if schedule.TermDescription == "Initial Payment" {
			continue
		}

		formatted := model.PaymentScheduleViewModel{
			ID:                      schedule.ID,
			AssessmentNumber:        schedule.AssessmentNumber,
			TermDescription:         schedule.TermDescription,
			DueDate:                 schedule.DueDate,
			ExpectedAmount:          schedule.ExpectedAmount,
			ExpectedAmountFormatted: formatAmount(schedule.ExpectedAmount),
			SortOrder:               schedule.SortOrder,
		}
		log.Printf("Formatted payment schedule: %+v", formatted)
		formattedSchedules = append(formattedSchedules, formatted)
	}

	// Render the template with formatted data
	return ctx.Render("partials/student_info", fiber.Map{
		"Student":          studentInfo.Student,
		"YearLevel":        studentInfo.YearLevel,
		"GradesSummary":    studentInfo.GradesSummary,
		"Assessment":       formattedAssessment,
		"PaymentSchedules": formattedSchedules,
	})
}
