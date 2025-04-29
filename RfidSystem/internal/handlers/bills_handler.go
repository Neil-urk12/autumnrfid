package handlers

import (
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

// billsCache is an LRU cache for storing student bills data.
// It has a capacity of 5 items and a time-to-live of 1 hour.
var billsCache = NewLRUCache(5, time.Hour)

// formatAmount formats a float64 amount to a string with two decimal places.
func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// formatNullableAmount formats a nullable float64 amount to a string with two decimal places. It returns "0.00" if the amount is nil.
func formatNullableAmount(amount *float64) string {
	if amount == nil {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", *amount)
}

// formatAssessmentForView formats a model.Assessment into a model.AssessmentViewModel
// for display purposes, handling nil assessments and formatting amounts.
func formatAssessmentForView(assessment *model.Assessment) model.AssessmentViewModel {
	// Handle nil assessment to prevent nil pointer dereference
	if assessment == nil {
		return model.AssessmentViewModel{
			TotalFeeAmount:      "0.00",
			NetAssessmentAmount: "0.00",
			InitialPayment:      "0.00",
			TotalPaymentAmount:  "0.00",
			RemainingBalance:    "0.00",
			TotalDiscountAmount: "0.00",
			FullPmtIfB4Prelim:   "0.00",
			PerExamFee:          "0.00",
		}
	}

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

// HandleBills handles HTTP requests to retrieve and display student bills.
// It first checks the cache, then fetches data from the repository if not found.
// It expects a student ID via form value "rfid" or query parameter "student-id".
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

			// Check if Assessment is nil before formatting
			if billsData.Assessment == nil {
				log.Printf("[WARNING] Cached bills data for %s has nil Assessment", studentId)
				// Remove invalid cache entry
				billsCache.Delete(studentId)
			} else {
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
	}

	// update to follow this format C23-2Number-4Number-MAN121
	// match, err := regexp.MatchString(`^ACLC-\d{4}-\d{3}$`, studentId)
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).SendString("Error validating student ID")
	// }

	// if !match {
	// 	return ctx.Status(fiber.StatusBadRequest).SendString("Invalid student ID format. Must be in format ACLC-YYYY-XXX")
	// }

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

	// Check if Assessment is nil
	if billsData.Assessment == nil {
		log.Printf("Assessment is nil for student ID: %s\n", studentId)
		return ctx.Status(fiber.StatusNotFound).SendString("No assessment data found for this student")
	}

	// Store in cache
	billsCache.Set(studentId, billsData)

	log.Printf("Successfully retrieved bills data for student ID: %s\n", studentId)
	fmt.Println("why", billsData)

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
