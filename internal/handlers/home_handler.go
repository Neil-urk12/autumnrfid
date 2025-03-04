package handlers

import "github.com/gofiber/fiber/v2"

type AppHandler struct{}

func NewHandler() *AppHandler {
	return &AppHandler{}
}

func (app *AppHandler) HandleGetIndex(ctx *fiber.Ctx) error {
	context := fiber.Map{}
	return ctx.Render("home", context)
}
