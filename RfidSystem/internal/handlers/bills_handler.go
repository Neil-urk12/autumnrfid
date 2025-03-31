package handlers

import (
	"fmt"
	"regexp"
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleBills(ctx *fiber.Ctx) error {
	studentId := ctx.Query("student-id")

	if studentId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student Id is required")
	}

	match, err := regexp.MatchString(`^ACLC-\d{4}-\d{3}$`, studentId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error validating student ID")
	}

	if !match {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid student ID format. Must be in format ACLC-YYYY-XXX")
	}

	fmt.Printf("Starting bills handler for student ID: %s\n", studentId)

	billsRepo := repositories.NewRFIDRepository(h.db)
	fmt.Printf("Created bills repository\n")

	billsData, err := billsRepo.GetStudentBillsByRFID(studentId)
	fmt.Printf("GetStudentBillsByRFID result - err: %v, billsData: %+v\n", err, billsData)

	if err != nil {
		fmt.Printf("Error getting bills data: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal server error: %v", err))
	}

	if billsData == nil {
		fmt.Printf("No bills data found for student ID: %s\n", studentId)
		return ctx.Status(fiber.StatusNotFound).SendString("No bills data found for this student")
	}

	fmt.Printf("Successfully retrieved bills data for student ID: %s\n", studentId)

	err = ctx.Render("partials/bills", fiber.Map{
		"Title": "Student Bills",
		"Bills": billsData,
	})
	if err != nil {
		fmt.Printf("Template rendering error: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Template error: %v", err))
	}
	return nil
}
