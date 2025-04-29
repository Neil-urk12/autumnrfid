package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/model"
)

// DEPRECATED

// GetAllStudents retrieves a paginated list of all students from the database.
// It returns a slice of Student pointers and an error if the query fails.
func (r *RFIDRepository) GetAllStudents(page int) ([]*model.Student, error) {
	var students []*model.Student
	limit := 5
	offset := (page - 1) * limit

	query := `
		SELECT student_ID, department_ID, first_Name, last_Name, middle_Name, birthday, contact_number, email, year_Level, program, block_section, first_access_timestamp, last_access_timestamp
		FROM Students ORDER BY last_Name ASC LIMIT ? OFFSET ?;
	`

	log.Print(query)

	rows, err := r.dbClient.DB.Query(query, limit, offset)
	if err != nil {
		log.Printf("Error querying students: %v\n", err)
		return nil, fmt.Errorf("error querying students: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		student := &model.Student{}
		err := rows.Scan(
			&student.StudentID,
			&student.DepartmentID,
			&student.FirstName,
			&student.LastName,
			&student.MiddleName,
			&student.Birthday,
			&student.ContactNumber,
			&student.Email,
			&student.YearLevel,
			&student.Program,
			&student.BlockSection,
			&student.FirstAccessTimestamp,
			&student.LastAccessTimestamp,
		)
		if err != nil {
			log.Printf("Error scanning student row: %v\n", err)
			return nil, fmt.Errorf("error scanning student row: %v", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after scanning rows: %v\n", err)
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	if len(students) == 0 {
		log.Println("No students found in the database")
		return nil, nil
	}

	return students, nil
}

// GetStudentsForAssessmentTerm retrieves a paginated list of students for a specific assessment term.
// It returns a slice of StudentAssessmentSummary pointers, the total count of students for the term, and an error.
//
// Deprecated
func (r *RFIDRepository) GetStudentsForAssessmentTerm(termID int64, page, limit int) ([]*model.StudentAssessmentSummary, int, error) {
	var students []*model.StudentAssessmentSummary
	var totalStudents int

	// Calculate offset
	offset := (page - 1) * limit

	// Query to get the total count of students for the term
	countQuery := `SELECT COUNT(DISTINCT s.student_ID) FROM Students s JOIN Assessment a ON s.student_ID = a.student_ID WHERE a.term_id = ?`
	err := r.dbClient.DB.QueryRow(countQuery, termID).Scan(&totalStudents)
	if err != nil {
		log.Printf("Error querying total student count for assessment term %d: %v\n", termID, err)
		return nil, 0, fmt.Errorf("error querying total student count: %v", err)
	}

	// If no students, return early
	if totalStudents == 0 {
		log.Printf("No students found for assessment term %d\n", termID)
		return []*model.StudentAssessmentSummary{}, 0, nil // Return empty slice and 0 count
	}

	// Query to get the paginated list of students
	query := `
		SELECT
			s.student_ID AS StudentID,
			CONCAT(s.first_Name, ' ', s.last_Name) AS Name,
			s.program AS Course,
			CASE s.year_Level
				WHEN 1 THEN '1st Year'
				WHEN 2 THEN '2nd Year'
				WHEN 3 THEN '3rd Year'
				WHEN 4 THEN '4th Year'
				ELSE CONCAT(s.year_Level, 'th Year')
			END AS YearLevel,
			s.status AS status
		FROM
			Students s
		JOIN
			Assessment a ON s.student_ID = a.student_ID
		WHERE
			a.term_id = ?
		ORDER BY
			s.last_Name, s.first_Name
		LIMIT ? OFFSET ?;
	`

	log.Printf("Executing query: %s with term_id: %d, limit: %d, offset: %d\n", query, termID, limit, offset)

	rows, err := r.dbClient.DB.Query(query, termID, limit, offset)
	if err != nil {
		log.Printf("Error querying students for assessment term %d: %v\n", termID, err)
		return nil, 0, fmt.Errorf("error querying students for assessment term: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		student := &model.StudentAssessmentSummary{}
		err := rows.Scan(
			&student.StudentID,
			&student.Name,
			&student.Course,
			&student.YearLevel,
			&student.Status,
		)
		if err != nil {
			log.Printf("Error scanning student assessment summary row: %v\n", err)
			return nil, 0, fmt.Errorf("error scanning student assessment summary row: %v", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after scanning student assessment summary rows: %v\n", err)
		return nil, 0, fmt.Errorf("error after scanning student assessment summary rows: %v", err)
	}

	// Note: We already checked totalStudents earlier, so len(students) might be 0 for later pages, which is expected.

	return students, totalStudents, nil
}

// Deprecated
func (r *RFIDRepository) GetStudentGradesByID(studentId string) (*model.GradesRecord, error) {
	// grades := &model.GradesRecord{}

	query := "CALL GetStudent(?)"

	rows, err := r.dbClient.DB.Query(query, studentId)
	if err != nil {
		log.Printf("Error querying student grades: %v\n", err)
		return nil, fmt.Errorf("error querying student grades: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

	}

	return nil, nil
}

// GetStudentInfo retrieves detailed information for a specific student by their ID.
// It queries the Students table and returns a Student pointer or an error.
// It returns sql.ErrNoRows if no student is found.
func (r *RFIDRepository) GetStudentInfo(studentID string) (*model.Student, error) {
	student := &model.Student{}

	query := `
		SELECT
			student_ID, department_ID, first_Name, last_Name, middle_Name,
			birthday, contact_number, email, year_Level, program, block_section,
			first_access_timestamp, last_access_timestamp, status
		FROM Students
		WHERE student_ID = ?;
	`

	log.Printf("Executing query: %s with student_id: %s\n", query, studentID) // Log string ID

	row := r.dbClient.DB.QueryRow(query, studentID) // Pass string ID
	err := row.Scan(
		&student.StudentID,
		&student.DepartmentID,
		&student.FirstName,
		&student.LastName,
		&student.MiddleName,
		&student.Birthday,
		&student.ContactNumber,
		&student.Email,
		&student.YearLevel,
		&student.Program,
		&student.BlockSection,
		&student.FirstAccessTimestamp,
		&student.LastAccessTimestamp,
		&student.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No student found with ID: %s\n", studentID) // Log string ID
			return nil, fmt.Errorf("student not found")             // Return a specific error for not found
		}
		log.Printf("Error scanning student row for ID %s: %v\n", studentID, err) // Log string ID
		return nil, fmt.Errorf("error scanning student row: %v", err)
	}

	return student, nil
}
