package repositories

import (
	"database/sql"
	"fmt"
)

type Student struct {
	StudentID    string
	DepartmentID int
	FirstName    string
	LastName     string
	MiddleName   string
	YearLevel    int
	Program      string
}

type RFIDRepository struct {
	dbClient *DatabaseClient
}

func NewRFIDRepository(dbClient *DatabaseClient) *RFIDRepository {
	return &RFIDRepository{dbClient: dbClient}
}

func (r *RFIDRepository) GetStudentByRFID(studentId string) (*Student, error) {
	query := `
		SELECT student_ID, department_ID, first_Name, last_Name, middle_Name, year_Level, program
		FROM Students
		WHERE student_ID = ?
	`

	fmt.Printf("Executing query with student ID: %s\n", studentId)

	student := &Student{}
	err := r.dbClient.DB.QueryRow(query, studentId).Scan(
		&student.StudentID,
		&student.DepartmentID,
		&student.FirstName,
		&student.LastName,
		&student.MiddleName,
		&student.YearLevel,
		&student.Program,
	)

	if err == sql.ErrNoRows {
		fmt.Printf("No student found with ID: %s\n", studentId)
		return nil, nil
	}
	if err != nil {
		fmt.Printf("Error querying student: %v\n", err)
		return nil, fmt.Errorf("error querying student: %v", err)
	}

	return student, nil
}
