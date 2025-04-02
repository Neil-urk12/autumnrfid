package handlers

import (
	"bufio"
	"fmt"
	"rfidsystem/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	rfid := ctx.FormValue("rfid")
	if rfid == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil {
		GetBroadcaster().Broadcast("error", fmt.Sprintf(`{"message": "Database error: %v"}`, err))
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Database error: %v", err))
	}

	if student == nil {
		htmxInstruction := fmt.Sprintf(`<div hx-get="/error" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`)
		GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Student not found: %s", rfid))
	}

	fmt.Printf("Student found: %s\n", student.StudentID)

	htmxInstruction := fmt.Sprintf(`<div hx-get="/student-partial/%s" hx-trigger="load" hx-swap="innerHTML" hx-target="#main"></div>`, rfid)

	GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
	return ctx.SendString("Processing")
}

func (h *AppHandler) GetStudentPartial(ctx *fiber.Ctx) error {
	rfid := ctx.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return ctx.Render("error_page", fiber.Map{})
	}

	return ctx.Render("partials/student_info", fiber.Map{
		"Student": student,
	}, "")
}

func getYearLevelString(year int) string {
	switch year {
	case 1:
		return "First"
	case 2:
		return "Second"
	case 3:
		return "Third"
	case 4:
		return "Fourth"
	default:
		return fmt.Sprintf("%dth", year)
	}
}

func (h *AppHandler) HandleStudentPartial(c *fiber.Ctx) error {
	rfid := c.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return c.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	return c.Render("partials/student_info", fiber.Map{
		"Student": student,
	}, "")
}

func (h *AppHandler) HandleSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	broadcaster := GetBroadcaster()

	client := &Client{
		messages: make(chan Message, 20),
		done:     make(chan struct{}),
	}

	broadcaster.register <- client

	c.Context().SetUserValue("client", client)

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			broadcaster.unregister <- client
			close(client.done)
		}()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		initialMsg := formatSSEMessage("connected", fmt.Sprintf(`{"time": "%s", "status": "connected"}`,
			time.Now().Format(time.RFC3339)))

		if _, err := w.WriteString(initialMsg); err != nil {
			return
		}

		if err := w.Flush(); err != nil {
			return
		}

		for {
			select {
			case <-client.done:
				fmt.Println("SSE connection closing (done channel triggered)")
				return
			case <-ticker.C:
				pingMsg := formatSSEMessage("ping", fmt.Sprintf(`{"time": "%s"}`, time.Now().Format(time.RFC3339)))
				if _, err := w.WriteString(pingMsg); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
				fmt.Println("Ping message sent successfully")
			case msg, ok := <-client.messages:
				if !ok {
					fmt.Println("SSE connection closing (channel closed)")
					return
				}

				// Format and send message
				sseMsg := formatSSEMessage(msg.Event, msg.Data)
				if _, err := w.WriteString(sseMsg); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	})

	return nil
}

func formatSSEMessage(event, data string) string {
	return fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
}
