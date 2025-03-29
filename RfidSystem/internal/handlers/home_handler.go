package handlers

import (
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	db *repositories.DatabaseClient
}

func NewHandler(db *repositories.DatabaseClient) *AppHandler {
	return &AppHandler{db: db}
}

func (h *AppHandler) HandleGetIndex(ctx *fiber.Ctx) error {
	return ctx.Render("pages/home", fiber.Map{})
}
