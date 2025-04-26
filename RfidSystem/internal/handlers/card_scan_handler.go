package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"rfidsystem/internal/repositories"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var (
	cardScanCache = NewLRUCache(5, time.Hour)
	cacheMutex    = &sync.RWMutex{}
)

func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	var req struct {
		RFID string `json:"rfid" form:"rfid"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	rfid := req.RFID
	log.Println(rfid, "RfID")
	if rfid == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	// Try cache first
	cacheMutex.RLock()
	cached, found := cardScanCache.Get(rfid)
	cacheMutex.RUnlock()
	if found {
		// Type assertion to view model
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Cache type assertion failed: %v", r)
			}
		}()
		student, ok := cached.(*model.StudentInfoViewModel)
		if ok && student != nil {
			htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)
			GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
			return ctx.SendString("Processing (cache)")
		}
	}

	rfidRepo := repositories.NewRFIDRepository(h.db)
	// student, err := rfidRepo.GetStudentByRFID(rfid)
	student, err := rfidRepo.GetStudentSummaryData(rfid)

	if err != nil {
		GetBroadcaster().Broadcast("error", fmt.Sprintf(`{"message": "Database error: %v"}`, err))
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Database error: %v", err))
	}

	if student == nil {
		htmxInstruction := fmt.Sprintf(`<div hx-get="/error" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`)
		GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Student not found: %s", rfid))
	}

	// Store in cache
	cacheMutex.Lock()
	cardScanCache.Set(rfid, student)
	cacheMutex.Unlock()

	htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)

	GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
	return ctx.SendString("Processing")
}

// HandleCardScanWS handles websocket card scan requests with bidirectional communication.
func (h *AppHandler) HandleCardScanWS(c *websocket.Conn) {
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}
		// Parse JSON message
		var payload struct {
			Status string `json:"status"`
			CardId string `json:"cardId"`
		}
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Printf("WS JSON parse error: %v", err)
			c.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}
		// If absent, render home page
		if payload.Status == "absent" {
			htmxInstruction := `<div hx-get="/" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`
			GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
			c.WriteMessage(websocket.TextMessage, []byte(htmxInstruction))
			continue
		}
		// Use CardId as RFID
		rfid := payload.CardId
		if rfid == "" {
			c.WriteMessage(websocket.TextMessage, []byte("RFID is required"))
			continue
		}
		// Try cache first
		cacheMutex.RLock()
		cached, found := cardScanCache.Get(rfid)
		cacheMutex.RUnlock()
		if found {
			if s, ok := cached.(*model.StudentInfoViewModel); ok && s != nil {
				htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)
				c.WriteMessage(websocket.TextMessage, []byte(htmxInstruction))
				c.WriteMessage(websocket.TextMessage, []byte("Processing (cache)"))
				GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
				continue
			}
		}
		// Fetch from DB
		student, err := h.RFIDRepository.GetStudentSummaryData(rfid)
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Database error: %v", err)))
			continue
		}
		if student == nil {
			htmxInstruction := `<div hx-get="/error" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`
			c.WriteMessage(websocket.TextMessage, []byte(htmxInstruction))
			GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
			continue
		}
		// Store in cache
		cacheMutex.Lock()
		cardScanCache.Set(rfid, student)
		cacheMutex.Unlock()
		htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)
		GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
		c.WriteMessage(websocket.TextMessage, []byte(htmxInstruction))
		c.WriteMessage(websocket.TextMessage, []byte("Processing"))
	}
}
