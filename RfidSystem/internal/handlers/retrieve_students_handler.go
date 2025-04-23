package handlers

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var studentsPageCache = NewLRUCache(5, time.Hour)

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
func (h *AppHandler) RetrieveStudentsHandler(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page", 1)
	cacheKey := strconv.Itoa(page)
	if cached, found := studentsPageCache.Get(cacheKey); found {
		students, ok := cached.([]interface{})
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
