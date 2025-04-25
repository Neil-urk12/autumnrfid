package handlers

import (
	"log"
	"math"
	"net/http"
	"rfidsystem/internal/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const defaultPage = 1
const defaultLimit = 5 // Default number of items per page

// HandleGetStudentsForAssessmentTerm handles requests to get a paginated list of students for a specific assessment term.
func (h *AppHandler) HandleGetStudentsForAssessmentTerm(c *fiber.Ctx) error {
	// Get termID from path parameter
	termIDStr := c.Params("termID")
	if termIDStr == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing termID parameter")
	}
	termID, err := strconv.ParseInt(termIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid termID parameter: %v", err)
		return c.Status(http.StatusBadRequest).SendString("Invalid termID parameter")
	}

	// Get page and limit from query parameters, with defaults
	page, err := strconv.Atoi(c.Query("page", strconv.Itoa(defaultPage)))
	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(c.Query("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 {
		limit = defaultLimit
	}

	// Fetch students and total count from repository
	students, totalStudents, err := h.RFIDRepository.GetStudentsForAssessmentTerm(termID, page, limit)
	if err != nil {
		log.Printf("Error getting students for assessment term %d (page %d, limit %d): %v", termID, page, limit, err)
		return c.Status(http.StatusInternalServerError).SendString("Error retrieving students")
	}

	// Calculate pagination metadata
	totalPages := 0
	if totalStudents > 0 {
		totalPages = int(math.Ceil(float64(totalStudents) / float64(limit)))
	}

	// Ensure students is an empty slice, not nil, if no results
	if students == nil {
		students = []*model.StudentAssessmentSummary{}
	}

	// Construct the response
	response := model.PaginatedStudentAssessmentResponse{
		Data: students,
		Pagination: model.PaginationMetadata{
			CurrentPage: page,
			PageSize:    limit,
			TotalItems:  totalStudents,
			TotalPages:  totalPages,
		},
	}

	return c.JSON(response)
}
