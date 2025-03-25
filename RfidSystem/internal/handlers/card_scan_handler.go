package handlers

import (
	"bufio"
	"fmt"
	"rfidsystem/internal/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SSE broadcaster for sending messages to connected clients
// type Broadcaster struct {
// 	// clients tracks all connected SSE clients with a bool indicating active status
// 	clients map[*Client]bool
// 	// register receives new clients that want to subscribe to SSE events
// 	register chan *Client
// 	// unregister receives clients that should be removed from broadcasting
// 	unregister chan *Client
// 	// broadcast channel receives messages that should be sent to all connected clients
// 	broadcast chan Message
// 	// mutex protects concurrent access to the clients map
// 	mutex sync.RWMutex
// 	// done channel signals when the broadcaster should shut down
// 	done chan struct{}
// }

// type Message struct {
// 	Event string
// 	Data  string
// }

// type Client struct {
// 	// messages is a buffered channel that receives SSE messages for this specific client
// 	// Buffer size of 20 allows for message queuing before potential backpressure
// 	messages chan Message
// 	// done signals when this client's connection should be terminated
// 	// closed when the client disconnects or encounters an error
// 	done chan struct{}
// }

// var (
// 	broadcaster *Broadcaster
// 	once        sync.Once
// )

// func GetBroadcaster() *Broadcaster {
// 	once.Do(func() {
// 		broadcaster = &Broadcaster{
// 			clients:    make(map[*Client]bool),
// 			register:   make(chan *Client, 10),
// 			unregister: make(chan *Client, 10),
// 			broadcast:  make(chan Message, 100),
// 			done:       make(chan struct{}),
// 		}
// 		go broadcaster.run()
// 	})
// 	return broadcaster
// }

// func (b *Broadcaster) Close() {
// 	close(b.done)
// }

// func (b *Broadcaster) run() {
// 	ticker := time.NewTicker(30 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-b.done:
// 			for client := range b.clients {
// 				close(client.messages)
// 				delete(b.clients, client)
// 			}
// 			b.Close()
// 			return

// 		case client := <-b.register:
// 			b.mutex.Lock()
// 			b.clients[client] = true
// 			b.mutex.Unlock()
// 			fmt.Printf("Client registered, total clients: %d\n", len(b.clients))

// 		case client := <-b.unregister:
// 			b.mutex.Lock()
// 			if _, ok := b.clients[client]; ok {
// 				delete(b.clients, client)
// 				close(client.messages)
// 			}
// 			b.mutex.Unlock()
// 			fmt.Printf("Client unregistered, remaining clients: %d\n", len(b.clients))

// 		case message := <-b.broadcast:
// 			// Send to all clients concurrently
// 			b.mutex.RLock()
// 			for client := range b.clients {
// 				// Non-blocking send, skip clients with full buffers
// 				select {
// 				case client.messages <- message:
// 					fmt.Printf("Message sent successfully\n")
// 				default:
// 					// Client buffer full, consider unregistering
// 					go func(c *Client) {
// 						b.unregister <- c
// 					}(client)
// 				}
// 			}
// 			b.mutex.RUnlock()

// 		case <-ticker.C:
// 			// Periodic check for inactive clients and garbage collection
// 			// Clean up for leaked resources I think
// 			fmt.Printf("Active SSE clients: %d\n", len(b.clients))
// 		}
// 	}
// }

// func (b *Broadcaster) Broadcast(event string, data string) {
// 	if event == "" {
// 		event = "message"
// 	}

// 	message := Message{Event: event, Data: data}

// 	fmt.Printf("[SSE DEBUG] Raw message being sent:\n%s\n\n", message)

// 	select {
// 	case b.broadcast <- message:
// 		fmt.Printf("[SSE DEBUG] Message sent successfully\n\n")
// 	default:
// 		fmt.Println("Warning: SSE message buffer full, dropping message")
// 	}
// }

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
		GetBroadcaster().Broadcast("not-found", fmt.Sprintf(`{"rfid": "%s"}`, rfid))
		return ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Student not found: %s", rfid))
		// return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
		// return ctx.Render("error_page", fiber.Map{})
	}

	fmt.Printf("Student found: %s\n", student.StudentID)

	htmxInstruction := fmt.Sprintf(`<div hx-get="/student-partial/%s" hx-trigger="load" hx-swap="innerHTML" hx-target="#student-data-container"></div>`, rfid)

	GetBroadcaster().Broadcast("studentcallback", htmxInstruction)
	return ctx.SendString("Processing")
	// htmxInstruction := fmt.Sprintf("<div hx-get=\"/ui//fragments/%s.html\" hx-trigger=\"load\" hx-swap=\"innerHTML\" hx-target=\"#student-data-container\"></div>", fragmentToRender)

	// studentData := fmt.Sprintf(`{
	// 	"studentID":"%s",
	// 	"firstName":"%s",
	// 	"lastName":"%s",
	// 	"middleName":"%s",
	// 	"yearLevel":%d,
	// 	"yearLevelStr":"%s",
	// 	"program":"%s",
	// 	"grades":[
	// 		{"subject":"Mathematics", "code":"MATH101", "units":3, "grade":"1.00", "remarks":"Passed"},
	// 		{"subject":"English", "code":"ENG101", "units":3, "grade":"1.25", "remarks":"Passed"},
	// 		{"subject":"Science", "code":"SCI101", "units":4, "grade":"1.50", "remarks":"Passed"},
	// 		{"subject":"History", "code":"HIST101", "units":3, "grade":"1.75", "remarks":"Passed"},
	// 		{"subject":"Physical Education", "code":"PE101", "units":2, "grade":"1.00", "remarks":"Passed"}
	// 	],
	// 	"bills":[
	// 		{"description":"Tuition Fee", "amount":15000.00, "status":"Paid", "dueDate":"2024-01-15"},
	// 		{"description":"Library Fee", "amount":1500.00, "status":"Paid", "dueDate":"2024-01-15"},
	// 		{"description":"Laboratory Fee", "amount":2500.00, "status":"Unpaid", "dueDate":"2024-02-15"},
	// 		{"description":"Miscellaneous Fee", "amount":1000.00, "status":"Unpaid", "dueDate":"2024-02-15"}
	// 	]
	// }`,
	// 	student.StudentID,
	// 	student.FirstName,
	// 	student.LastName,
	// 	student.MiddleName,
	// 	student.YearLevel,
	// 	getYearLevelString(student.YearLevel),
	// 	student.Program)
}

func (h *AppHandler) GetStudentPartial(ctx *fiber.Ctx) error {
	rfid := ctx.Params("rfid")

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)

	if err != nil || student == nil {
		return ctx.Render("error_page", fiber.Map{})
	}

	return ctx.Render("partials/student_info", fiber.Map{
		"Student":   student,
		"YearLevel": getYearLevelString(student.YearLevel),
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
		"Student":   student,
		"YearLevel": getYearLevelString(student.YearLevel),
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
