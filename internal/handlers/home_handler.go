package handlers

import "github.com/gofiber/fiber/v2"

type HomeHandler struct{}

func NewHandler() *HomeHandler {
	return &HomeHandler{}
}

func (app *HomeHandler) HandleGetIndex(ctx *fiber.Ctx) error {
	context := fiber.Map{}
	return ctx.Render("home", context)
}
