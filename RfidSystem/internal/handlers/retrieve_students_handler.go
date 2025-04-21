package handlers

import (
	"log"
	"rfidsystem/internal/repositories"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var studentsPageCache = NewLRUCache(5, time.Hour)

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

	rfidRepo := repositories.NewRFIDRepository(h.db)
	students, err := rfidRepo.GetAllStudents(page)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	// Store in cache
	studentsPageCache.Set(cacheKey, students)

	return ctx.JSON(students)
}
