package repositories

import (
	"fmt"
	"log"
	"rfidsystem/internal/model"
)

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
