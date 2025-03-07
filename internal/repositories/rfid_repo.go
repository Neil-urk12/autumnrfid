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
		SELECT student_ID, DepartmentID, FirstName, LastName, MiddleName, YearLevel, Program
		FROM Students 
		WHERE student_ID = ?
	`

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
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying student: %v", err)
	}

	return student, nil
}
