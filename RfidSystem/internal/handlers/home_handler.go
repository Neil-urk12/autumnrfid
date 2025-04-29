package handlers

import (
	"rfidsystem/internal/repositories"

	"github.com/gofiber/fiber/v2"
)

type AppHandler struct {
	db             *repositories.DatabaseClient
	RFIDRepository *repositories.RFIDRepository
}

func NewHandler(db *repositories.DatabaseClient, rfidRepo *repositories.RFIDRepository) *AppHandler {
	return &AppHandler{db: db, RFIDRepository: rfidRepo}
}

func (h *AppHandler) HandleGetIndex(ctx *fiber.Ctx) error {
	return ctx.Render("pages/home", fiber.Map{})
}
