package handlers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
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

// HandleLogPartial renders only the log container for HTMX auto-refresh with filters
func (h *AppHandler) HandleLogPartial(c *fiber.Ctx) error {
	// Fetch with search and level filter
	search := strings.TrimSpace(c.Query("search", ""))
	level  := c.Query("level", "all")
	query := `SELECT id, timestamp, card_id, student_ID, event_type, message, details, status
       FROM scan_logs
       WHERE (? = '' OR message LIKE ? OR event_type LIKE ?)
         AND (? = 'all' OR status = ?)
       ORDER BY timestamp DESC`
	rows, err := h.db.DB.Query(query,
		search, "%"+search+"%", "%"+search+"%",
		level, level)
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

// HandleClearLogs transfers all scan_logs to archived_logs and deletes them
func (h *AppHandler) HandleClearLogs(c *fiber.Ctx) error {
	tx, err := h.db.DB.Begin()
	if err != nil {
		log.Printf("HandleClearLogs Begin error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to begin transaction: %v", err))
	}

	// Archive logs to archived_logs table
	_, err = tx.Exec(`INSERT INTO archived_logs (id, timestamp, card_id, student_ID, event_type, message, details, status)
       SELECT id, timestamp, card_id, student_ID, event_type, message, details, status FROM scan_logs`)
	if err != nil {
		tx.Rollback()
		log.Printf("HandleClearLogs archive error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to archive logs: %v", err))
	}

	// Delete original logs
	_, err = tx.Exec("DELETE FROM scan_logs")
	if err != nil {
		tx.Rollback()
		log.Printf("HandleClearLogs delete error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to delete logs: %v", err))
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Printf("HandleClearLogs commit error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to commit transaction: %v", err))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// HandleExportLogs returns all logs as CSV download
func (h *AppHandler) HandleExportLogs(c *fiber.Ctx) error {
	rows, err := h.db.DB.Query(`SELECT id, timestamp, card_id, student_ID, event_type, message, details, status
       FROM scan_logs ORDER BY timestamp DESC`)
	if err != nil {
		log.Printf("HandleExportLogs Query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).
			SendString(fmt.Sprintf("Failed to query logs for export: %v", err))
	}
	defer rows.Close()
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=logs.csv")
	writer := csv.NewWriter(c)
	defer writer.Flush()
	// header
	writer.Write([]string{"ID","Timestamp","CardID","StudentID","EventType","Message","Details","Status"})
	for rows.Next() {
		var id int
		var ts time.Time
		var cardID string
		var student sql.NullString
		var eventType, message string
		var details sql.NullString
		var status string
		if err := rows.Scan(&id, &ts, &cardID, &student, &eventType, &message, &details, &status); err != nil {
			log.Printf("HandleExportLogs ScanRow error: %v", err)
			continue
		}
		record := []string{
			strconv.Itoa(id),
			ts.Format(time.RFC3339),
			cardID,
			student.String,
			eventType,
			message,
			details.String,
			status,
		}
		writer.Write(record)
	}
	return nil
}
