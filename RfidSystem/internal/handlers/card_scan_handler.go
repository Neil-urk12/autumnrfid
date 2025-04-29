package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// LRU Cache for card scans
var cardScanCache = NewLRUCache(5, time.Hour)

// HandleCardScan handles HTTP POST requests for RFID card scans.
// It processes the RFID from the request body or form, logs the event,
// checks the cache, fetches student data from the repository if necessary,
// stores the data in the cache, and broadcasts an HTMX instruction via SSE.
func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	var req struct {
		RFID string `json:"rfid" form:"rfid"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		if err2 := h.db.LogScanEvent(req.RFID, nil, "card_read_error", fmt.Sprintf("Error parsing request body: %v", err), "", "failure"); err2 != nil {
			log.Printf("Failed to log event: %v", err2)
		}
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}
	rfid := req.RFID
	log.Println(rfid, "RfID")
	if rfid != "" {
		if err2 := h.db.LogScanEvent(rfid, nil, "scan", fmt.Sprintf("Card scanned: %s", rfid), "", "info"); err2 != nil {
			log.Printf("Failed to log scan event: %v", err2)
		}
	}
	if rfid == "" {
		_ = h.db.LogScanEvent("", nil, "card_read_error", "RFID is required", "", "failure")
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	// Try cache first
	cached, found := cardScanCache.Get(rfid)
	if found {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Cache type assertion failed: %v", r)
			}
		}()
		student, ok := cached.(*model.StudentInfoViewModel)
		if ok && student != nil {
			// Log cache hit event
			_ = h.db.LogScanEvent(rfid, &student.Student.StudentID, "scan_cache_hit", fmt.Sprintf("Cache hit for student %s", student.Student.StudentID), "", "info")
			htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)
			GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
			return ctx.SendString("Processing (cache)")
		}
	}

	student, err := h.RFIDRepository.GetStudentSummaryData(rfid)

	if err != nil {
		_ = h.db.LogScanEvent(rfid, nil, "db_error", fmt.Sprintf("Database error: %v", err), "", "failure")
		GetBroadcaster().Broadcast("error", fmt.Sprintf(`{"message": "Database error: %v"}`, err))
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Database error: %v", err))
	}

	if student == nil {
		_ = h.db.LogScanEvent(rfid, nil, "student_not_found", fmt.Sprintf("Student not found: %s", rfid), "", "failure")
		htmxInstruction := `<div hx-get="/error" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`
		GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Student not found: %s", rfid))
	}

	// Store in cache
	cardScanCache.Set(rfid, student)

	htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)

	GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
	_ = h.db.LogScanEvent(rfid, &student.Student.StudentID, "info_displayed", fmt.Sprintf("Displayed info for student ID %s", student.Student.StudentID), "", "success")
	return ctx.SendString("Processing")
}

// HandleCardScanWS handles websocket connections for real-time RFID card scans.
// It reads messages from the WebSocket, processes the card ID,
// checks the cache, fetches student data if needed, stores in cache,
// broadcasts an HTMX instruction via SSE, and sends the instruction back over the WebSocket.
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
		cached, found := cardScanCache.Get(rfid)
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
		cardScanCache.Set(rfid, student)
		htmxInstruction := fmt.Sprintf(`<div hx-post="/student-partial" hx-vals='{"rfid":"%s"}' hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)
		GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
		c.WriteMessage(websocket.TextMessage, []byte(htmxInstruction))
		c.WriteMessage(websocket.TextMessage, []byte("Processing"))
	}
}
