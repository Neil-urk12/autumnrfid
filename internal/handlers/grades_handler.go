package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleGrades(ctx *fiber.Ctx) error {
	err := ctx.Render("pages/grades", fiber.Map{
		"Title": "Student Grades",
	})
	if err != nil {
		fmt.Printf("Template rendering error: %v\n", err)
		return err
	}
	return nil
}

func (h *AppHandler) HandleTestGrades(ctx *fiber.Ctx) error {
	err := ctx.Render("partials/grades", fiber.Map{
		"Title": "Student Grades",
	})
	if err != nil {
		fmt.Printf("Template rendering error: %v\n", err)
		return err
	}
	return nil
}
