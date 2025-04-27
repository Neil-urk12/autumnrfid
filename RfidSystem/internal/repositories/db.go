package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseClient struct {
	DB *sql.DB
}

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

// Close closes the database connection
func (c *DatabaseClient) Close() error {
	return c.DB.Close()
}

// type Test struct {
// 	StudentId    string
// 	DepartmentID int
// 	FirstName    string
// 	LastName     string
// 	MiddleName   string
// 	YearLevel    int
// 	Program      string
// }

// func (c *DatabaseClient) GetStudentByRFID(studentId string) (*Test, error) {
// 	query := `
// 		SELECT *
// 		FROM Students
// 		WHERE student_ID = ?
// 	`

// 	student := &Test{}
// 	err := c.DB.QueryRow(query, studentId).Scan(
// 		&student.StudentId,
// 		&student.DepartmentID,
// 		&student.FirstName,
// 		&student.LastName,
// 		&student.MiddleName,
// 		&student.YearLevel,
// 		&student.Program,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying student: %v", err)
// 	}

// 	return student, nil
// }

// func (c *DatabaseClient) GetStudentGrades(studentID int) ([]Grade, error) {
// 	query := `
// 		SELECT
// 		FROM Enrollments
// 		WHERE student_id = ?
// 	`

// 	rows, err := c.DB.Query(query, studentID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error querying grades: %v", err)
// 	}
// 	defer rows.Close()

// 	var grades []Grade
// 	for rows.Next() {
// 		var g Grade
// 		if err := rows.Scan(&g.SubjectCode, &g.SubjectName, &g.Grade); err != nil {
// 			return nil, fmt.Errorf("error scanning grade row: %v", err)
// 		}
// 		grades = append(grades, g)
// 	}

// 	return grades, nil
// }

// LogScanEvent inserts a new log entry into the scan_logs table.
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
