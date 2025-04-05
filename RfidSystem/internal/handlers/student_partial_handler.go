package handlers

import (
	"log"
	"rfidsystem/internal/model"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleStudentInfo(ctx *fiber.Ctx) error {
	studentId := ctx.Params("rfid")
	if studentId == "" {
		log.Printf("Student ID is required")
		return ctx.Status(fiber.StatusBadRequest).SendString("Student ID is required")
	}

	log.Printf("Retrieving data")
	rfidRepo := repositories.NewRFIDRepository(h.db)
	studentInfo, err := rfidRepo.GetStudentSummaryData(studentId)
	if err != nil {
		log.Printf("Error getting student info: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to retrieve student information")
	}

	if studentInfo == nil || studentInfo.Student == nil {
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}
	log.Printf("Student info retrieved successfully")

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
	formattedAssessment := formatAssessmentForView(studentInfo.Assessment)
	log.Printf("Formatted assessment: %+v", formattedAssessment)

	// Format payment schedules
	var formattedSchedules []model.PaymentScheduleViewModel
	for _, schedule := range studentInfo.PaymentSchedules {
		if schedule.TermDescription == "Initial" {
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

func (h *AppHandler) GetStudentPartial(ctx *fiber.Ctx) error {
	rfid := ctx.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return ctx.Render("error_page", fiber.Map{})
	}

	return ctx.Render("partials/student_info", fiber.Map{
		"Student": student,
	}, "")
}

func (h *AppHandler) HandleStudentPartial(c *fiber.Ctx) error {
	rfid := c.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return c.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	return c.Render("partials/student_info", fiber.Map{
		"Student": student,
	}, "")
}
