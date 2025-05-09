package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Deprecated

var studentsPageCache = NewLRUCache(5, time.Hour)

// GetStudentById handles HTTP requests to retrieve detailed information for a specific student by their ID.
// It expects the student ID as a path parameter.
func (h *AppHandler) GetStudentById(ctx *fiber.Ctx) error {
	studentID := ctx.Params("id")
	log.Printf("Received request to get student with ID: %s\n", studentID)

	student, err := h.RFIDRepository.GetStudentInfo(studentID)
	if err != nil {
		if err.Error() == "student not found" {
			log.Printf("Student not found with ID: %s\n", studentID)
			return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
		}
		log.Printf("Error retrieving student info for ID %s: %v\n", studentID, err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	if student == nil {
		log.Printf("Student not found with ID: %s (nil returned without error)\n", studentID)
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	return ctx.JSON(student)
}

// GetGrades handles HTTP requests to retrieve grades for a specific student by their ID.
// It expects the student ID as a path parameter.
func (h *AppHandler) GetGrades(ctx *fiber.Ctx) error {
	studentID := ctx.Params("id")
	log.Printf("Received request to get grades for student with ID: %s\n", studentID)

	grades, err := h.RFIDRepository.GetStudentGradesByID(studentID)
	if err != nil {
		if err.Error() == "student not found" {
			log.Printf("Student not found with ID: %s\n", studentID)
			return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
		}
		log.Printf("Error retrieving grades for student ID %s: %v\n", studentID, err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	if grades == nil {
		log.Printf("Grades not found for student ID: %s (nil returned without error)\n", studentID)
		return ctx.Status(fiber.StatusNotFound).SendString("Grades not found")
	}

	log.Println(ctx.JSON(grades))
	return ctx.JSON(grades)
}

// RetrieveStudentsHandler handles HTTP requests to retrieve a paginated list of students.
// It expects the page number as a query parameter.
func (h *AppHandler) RetrieveStudentsHandler(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	cacheKey := strconv.Itoa(page)
	if cached, found := studentsPageCache.Get(cacheKey); found {
		students, ok := cached.([]any)
		if ok && students != nil {
			log.Printf("[CACHE HIT] Students page %d", page)
			return ctx.JSON(students)
		}
	}

	students, err := h.RFIDRepository.GetAllStudents(page)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	// Store in cache
	studentsPageCache.Set(cacheKey, students)

	return ctx.JSON(students)
}
