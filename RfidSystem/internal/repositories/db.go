package repositories

import (
	"database/sql"
	"fmt"
	"rfidsystem/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseClient struct {
	DB *sql.DB
}

func NewDatabaseClient(config config.DatabaseConfig) (*DatabaseClient, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
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
