package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ScanLog struct {
	ID        int
	Timestamp time.Time
	CardID    string
	StudentID sql.NullString
	EventType string
	Message   string
	Details   sql.NullString
	Status    string
}

// HandleLog renders the LogMonitor page
func (h *AppHandler) HandleLog(c *fiber.Ctx) error {
	// Fetch logs from database
	rows, err := h.db.DB.Query(
		`SELECT id, timestamp, card_id, student_ID, event_type, message, details, status
		 FROM scan_logs
		 ORDER BY timestamp DESC`)
	if err != nil {
		log.Printf("HandleLog Query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to query logs: %v", err))
	}
	defer rows.Close()

	var logs []ScanLog
	for rows.Next() {
		var sl ScanLog
		if err := rows.Scan(&sl.ID, &sl.Timestamp, &sl.CardID, &sl.StudentID, &sl.EventType, &sl.Message, &sl.Details, &sl.Status); err != nil {
			log.Printf("HandleLog ScanRow error: %v", err)
			return c.Status(fiber.StatusInternalServerError).
				SendString(fmt.Sprintf("Failed to scan log row: %v", err))
		}
		logs = append(logs, sl)
	}
	if err := rows.Err(); err != nil {
		log.Printf("HandleLog rows iteration error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Error reading logs: %v", err))
	}

	total := len(logs)
	errorCount := 0
	warnCount := 0
	for _, sl := range logs {
		switch sl.Status {
		case "error":
			errorCount++
		case "warn", "warning":
			warnCount++
		}
	}
	var rate float64
	if total > 1 {
		duration := logs[0].Timestamp.Sub(logs[total-1].Timestamp).Seconds()
		if duration > 0 {
			rate = float64(total) / duration
		}
	}
	if err := c.Render("pages/log", fiber.Map{
		"Logs":      logs,
		"TotalLogs": total,
		"ErrorLogs": errorCount,
		"WarnLogs":  warnCount,
		"LogRate":   rate,
	}); err != nil {
		log.Printf("HandleLog Render error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Render error: %v", err))
	}
	return nil
}

// HandleLogPartial renders only the log container for HTMX auto-refresh
func (h *AppHandler) HandleLogPartial(c *fiber.Ctx) error {
	rows, err := h.db.DB.Query(
		`SELECT id, timestamp, card_id, student_ID, event_type, message, details, status
		 FROM scan_logs
		 ORDER BY timestamp DESC`)
	if err != nil {
		log.Printf("HandleLogPartial Query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to query logs: %v", err))
	}
	defer rows.Close()

	var logs []ScanLog
	for rows.Next() {
		var sl ScanLog
		if err := rows.Scan(&sl.ID, &sl.Timestamp, &sl.CardID, &sl.StudentID, &sl.EventType, &sl.Message, &sl.Details, &sl.Status); err != nil {
			log.Printf("HandleLogPartial ScanRow error: %v", err)
			return c.Status(fiber.StatusInternalServerError).
				SendString(fmt.Sprintf("Failed to scan log row: %v", err))
		}
		logs = append(logs, sl)
	}
	if err := rows.Err(); err != nil {
		log.Printf("HandleLogPartial rows error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Error reading logs: %v", err))
	}
	return c.Render("partials/log_list", fiber.Map{"Logs": logs})
}

// HandleStatsPartial renders only the stats cards for HTMX auto-refresh
func (h *AppHandler) HandleStatsPartial(c *fiber.Ctx) error {
	// Fetch logs from database
	rows, err := h.db.DB.Query(
		`SELECT timestamp, status FROM scan_logs ORDER BY timestamp DESC`)
	if err != nil {
		log.Printf("HandleStatsPartial Query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to query logs for stats: %v", err))
	}
	defer rows.Close()
	type entry struct {
		Timestamp time.Time
		Status    string
	}
	var entries []entry
	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.Timestamp, &e.Status); err != nil {
			log.Printf("HandleStatsPartial ScanRow error: %v", err)
			return c.Status(fiber.StatusInternalServerError).
				SendString(fmt.Sprintf("Failed to scan stat row: %v", err))
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		log.Printf("HandleStatsPartial rows err: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Error reading stats rows: %v", err))
	}
	// Compute stats
	total := len(entries)
	errors := 0
	warns := 0
	for _, e := range entries {
		switch e.Status {
		case "error":
			errors++
		case "warn", "warning":
			warns++
		}
	}
	// Compute rate logs/sec
	var rate float64
	if total > 1 {
		dt := entries[0].Timestamp.Sub(entries[total-1].Timestamp).Seconds()
		if dt > 0 {
			rate = float64(total) / dt
		}
	}
	// Render stats partial
	return c.Render("partials/stats", fiber.Map{
		"TotalLogs": total,
		"ErrorLogs": errors,
		"WarnLogs":  warns,
		"LogRate":   rate,
	})
}
