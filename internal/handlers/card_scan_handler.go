package handlers

import (
	"bufio"
	"fmt"
	"rfidsystem/internal/repositories"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SSE broadcaster for sending messages to connected clients
type Broadcaster struct {
	events chan string
	mu     sync.Mutex
}

// Global broadcaster instance
var (
	broadcaster *Broadcaster
	once        sync.Once
)

// GetBroadcaster returns the singleton broadcaster instance
func GetBroadcaster() *Broadcaster {
	once.Do(func() {
		broadcaster = &Broadcaster{
			events: make(chan string, 100),
		}
	})
	return broadcaster
}

// Broadcast sends a message with the specified event type to all SSE clients
func (b *Broadcaster) Broadcast(event string, data string) {
	// Make sure to follow SSE format EXACTLY: "event: name\ndata: data\n\n"
	// The double newline at the end is CRITICAL
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)

	// Debug what we're actually sending
	fmt.Printf("[SSE DEBUG] Raw message being sent:\n%s\n", message)

	// Immediately write to channel if there's room, otherwise drop the message
	select {
	case b.events <- message:
		// Message sent successfully
	default:
		fmt.Println("Warning: SSE message buffer full, dropping message")
	}
}

func (h *AppHandler) HandleCardScan(ctx *fiber.Ctx) error {
	rfid := ctx.FormValue("rfid")
	if rfid == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("RFID is required")
	}

	rfidRepo := repositories.NewRFIDRepository(h.db)
	student, err := rfidRepo.GetStudentByRFID(rfid)
	if err != nil {
		// Log the actual error
		fmt.Printf("Database error: %v\n", err)
		GetBroadcaster().Broadcast("error", fmt.Sprintf(`{"message": "Database error: %v"}`, err))
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Database error: %v", err))
	}

	// Handle not found case
	if student == nil {
		// Send a not-found event to all SSE clients
		GetBroadcaster().Broadcast("not-found", fmt.Sprintf(`{"rfid": "%s"}`, rfid))
		return ctx.Status(fiber.StatusNotFound).SendString("Student not found")
	}

	// Format student data as clean JSON for SSE - using compact format to avoid whitespace issues
	// Include dummy grades and bills data
	studentData := fmt.Sprintf(`{
		"studentID":"%s",
		"firstName":"%s",
		"lastName":"%s",
		"middleName":"%s",
		"yearLevel":%d,
		"yearLevelStr":"%s",
		"program":"%s",
		"grades":[
			{"subject":"Mathematics", "code":"MATH101", "units":3, "grade":"1.00", "remarks":"Passed"},
			{"subject":"English", "code":"ENG101", "units":3, "grade":"1.25", "remarks":"Passed"},
			{"subject":"Science", "code":"SCI101", "units":4, "grade":"1.50", "remarks":"Passed"},
			{"subject":"History", "code":"HIST101", "units":3, "grade":"1.75", "remarks":"Passed"},
			{"subject":"Physical Education", "code":"PE101", "units":2, "grade":"1.00", "remarks":"Passed"}
		],
		"bills":[
			{"description":"Tuition Fee", "amount":15000.00, "status":"Paid", "dueDate":"2024-01-15"},
			{"description":"Library Fee", "amount":1500.00, "status":"Paid", "dueDate":"2024-01-15"},
			{"description":"Laboratory Fee", "amount":2500.00, "status":"Unpaid", "dueDate":"2024-02-15"},
			{"description":"Miscellaneous Fee", "amount":1000.00, "status":"Unpaid", "dueDate":"2024-02-15"}
		]
	}`,
		student.StudentID,
		student.FirstName,
		student.LastName,
		student.MiddleName,
		student.YearLevel,
		getYearLevelString(student.YearLevel),
		student.Program)

	// Remove all whitespace to ensure proper JSON formatting for SSE
	studentData = removeWhitespace(studentData)

	// Send student data to all connected SSE clients
	GetBroadcaster().Broadcast("student-data", studentData)

	// Log that we've broadcast the event (for debugging)
	fmt.Printf("Broadcast sent with event: 'student-data' and data: %s\n", studentData)

	// Handle response based on content type
	if ctx.Accepts("application/json") == "application/json" {
		// For API clients, return JSON response
		fmt.Println("Returning JSON response")
		return ctx.JSON(fiber.Map{
			"status": "success",
			"data":   student,
		})
	} else {
		// For direct requests, render the template
		fmt.Println("Rendering template")
		return ctx.Render("home", fiber.Map{
			"Student":   student,
			"Title":     "Student Information",
			"YearLevel": getYearLevelString(student.YearLevel),
		})
	}
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

// Helper function to remove all whitespace from a string for compact JSON
func removeWhitespace(s string) string {
	// First remove all newlines and tabs
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	
	// Handle spaces more carefully - don't remove spaces in quoted strings
	var result strings.Builder
	inQuotes := false
	
	for _, r := range s {
		if r == '"' {
			inQuotes = !inQuotes
			result.WriteRune(r)
		} else if inQuotes || r != ' ' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// HandleSSE establishes a server-sent events connection
func (h *AppHandler) HandleSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Log SSE connection
	fmt.Printf("SSE connection established from %s\n", c.IP())

	// Send initial connection message
	initMessage := "event: connected\ndata: {\"time\": \"" + time.Now().Format(time.RFC3339) + "\", \"status\": \"connected\"}\n\n"

	// Get broadcaster instance
	broadcaster := GetBroadcaster()

	// Create channel for client-specific cleanup
	done := make(chan bool)

	// Setup cleanup when connection is closed
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		// Send initial message first
		if fw, err := w.Write([]byte(initMessage)); err != nil || fw == 0 {
			fmt.Printf("Error sending initial SSE message: %v\n", err)
			return
		}
		if err := w.Flush(); err != nil {
			fmt.Printf("Error flushing initial SSE message: %v\n", err)
			return
		}

		fmt.Println("Initial SSE message sent successfully")
		eventsChannel := broadcaster.events

		for {
			select {
			case <-done:
				fmt.Println("SSE connection closing (done channel triggered)")
				return
			case <-ticker.C:
				// Send heartbeat
				pingMsg := "event: ping\ndata: {\"time\": \"" + time.Now().Format(time.RFC3339) + "\"}\n\n"
				fw, err := w.Write([]byte(pingMsg))
				if err != nil || fw == 0 {
					fmt.Printf("Error sending SSE ping: %v\n", err)
					close(done)
					return
				}
				if err = w.Flush(); err != nil {
					fmt.Printf("Error flushing SSE ping: %v\n", err)
					close(done)
					return
				}
			case msg := <-eventsChannel:
				// Send message from broadcaster
				fmt.Printf("Broadcasting message: %s\n", msg)
				fw, err := w.Write([]byte(msg))
				if err != nil || fw == 0 {
					fmt.Printf("Error sending SSE message: %v\n", err)
					close(done)
					return
				}
				if err = w.Flush(); err != nil {
					fmt.Printf("Error flushing SSE message: %v\n", err)
					close(done)
					return
				}
				fmt.Println("SSE message sent successfully")
			}
		}
	})

	return nil
}
