package handlers

import (
	"fmt"
	"log"
	"regexp"
	"rfidsystem/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

var billsCache = NewLRUCache(5, time.Hour)

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func formatNullableAmount(amount *float64) string {
	if amount == nil {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", *amount)
}

func formatAssessmentForView(assessment *model.Assessment) model.AssessmentViewModel {
	return model.AssessmentViewModel{
		ID:                  assessment.ID,
		StudentID:           assessment.StudentID,
		TermID:              assessment.TermID,
		TotalFeeAmount:      formatAmount(assessment.TotalFeeAmount),
		NetAssessmentAmount: formatAmount(assessment.NetAssessmentAmount),
		InitialPayment:      formatNullableAmount(assessment.InitialPayment),
		TotalPaymentAmount:  formatAmount(assessment.TotalPaymentAmount),
		RemainingBalance:    formatAmount(assessment.RemainingBalance),
		TotalDiscountAmount: formatAmount(assessment.TotalDiscountAmount),
		FullPmtIfB4Prelim:   formatNullableAmount(assessment.FullPmtIfB4Prelim),
		PerExamFee:          formatNullableAmount(assessment.PerExamFee),
	}
}

func (h *AppHandler) HandleBills(ctx *fiber.Ctx) error {
	studentId := ctx.FormValue("rfid")
	if studentId == "" {
		studentId = ctx.Query("student-id")
	}
	if studentId == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Student Id is required")
	}
	// Try cache first
	if cached, found := billsCache.Get(studentId); found {
		billsData, ok := cached.(*model.Bills)
		if ok && billsData != nil {
			log.Printf("[CACHE HIT] Bills for %s", studentId)
			assessmentMap := formatAssessmentForView(billsData.Assessment)
			return ctx.Render("partials/bills", fiber.Map{
				"Title":          "Student Bills",
				"Bills":          assessmentMap,
				"FeeBreakdown":   billsData.FeeBreakdown,
				"Discounts":      billsData.Discounts,
				"PaymentHistory": billsData.PaymentHistory,
			})
		}
	}

	match, err := regexp.MatchString(`^ACLC-\d{4}-\d{3}$`, studentId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error validating student ID")
	}

	if !match {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid student ID format. Must be in format ACLC-YYYY-XXX")
	}

	log.Printf("Created bills repository\n")

	billsData, err := h.RFIDRepository.GetStudentBillsByRFID(studentId)
	log.Printf("GetStudentBillsByRFID result - err: %v, billsData: %+v\n", err, billsData)

	if err != nil {
		log.Printf("Error getting bills data: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal server error: %v", err))
	}

	if billsData == nil {
		log.Printf("No bills data found for student ID: %s\n", studentId)
		return ctx.Status(fiber.StatusNotFound).SendString("No bills data found for this student")
	}
	// Store in cache
	billsCache.Set(studentId, billsData)

	log.Printf("Successfully retrieved bills data for student ID: %s\n", studentId)

	assessmentMap := formatAssessmentForView(billsData.Assessment)

	err = ctx.Render("partials/bills", fiber.Map{
		"Title":          "Student Bills",
		"Bills":          assessmentMap,
		"FeeBreakdown":   billsData.FeeBreakdown,
		"Discounts":      billsData.Discounts,
		"PaymentHistory": billsData.PaymentHistory,
	})
	if err != nil {
		log.Printf("Template rendering error: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Template error: %v", err))
	}
	return nil
}
