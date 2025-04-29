package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseClient represents a client for interacting with the database.
type DatabaseClient struct {
	DB *sql.DB
}

// NewDatabaseClient creates a new DatabaseClient and establishes a database connection
// using the provided database configuration. It returns a pointer to the client and an error.
func NewDatabaseClient(config config.DatabaseConfig) (*DatabaseClient, error) {
	// Enable parsing of MySQL TIMESTAMP fields into time.Time
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.DatabaseName,
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DatabaseClient{DB: db}, nil
}

// Close closes the database connection held by the DatabaseClient.
func (c *DatabaseClient) Close() error {
	return c.DB.Close()
}

// LogScanEvent inserts a new log entry into the scan_logs table.
// It records details about a scan event, including the card ID, optional student ID,
// event type, message, optional details, and status.
func (c *DatabaseClient) LogScanEvent(cardID string, studentID *string, eventType, message, details, status string) error {
	log.Printf("[LogScanEvent] called with cardID=%s studentID=%v eventType=%s message=%q status=%s", cardID, studentID, eventType, message, status)
	ts := time.Now().UTC()
	var detailsParam interface{}
	if details == "" {
		detailsParam = nil
	} else {
		detailsParam = details
	}
	res, err := c.DB.Exec(
		`INSERT INTO scan_logs (timestamp, card_id, student_ID, event_type, message, details, status)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ts, cardID, studentID, eventType, message, detailsParam, status,
	)
	if err != nil {
		log.Printf("LogScanEvent SQL error: %v", err)
		return err
	}
	if affected, err2 := res.RowsAffected(); err2 != nil {
		log.Printf("LogScanEvent RowsAffected error: %v", err2)
	} else {
		log.Printf("LogScanEvent inserted rows: %d", affected)
	}
	return nil
}
