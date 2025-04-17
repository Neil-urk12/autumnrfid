package handlers

import (
	"bufio"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

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
				log.Println("SSE connection closing (done channel triggered)")
				return
			case <-ticker.C:
				pingMsg := formatSSEMessage("ping", fmt.Sprintf(`{"time": "%s"}`, time.Now().Format(time.RFC3339)))
				if _, err := w.WriteString(pingMsg); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					log.Printf("Failed to flush ping message: %v", err)
					return
				}
				log.Println("Ping message sent successfully")
			case msg, ok := <-client.messages:
				if !ok {
					log.Println("SSE connection closing (channel closed)")
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
