package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"rfidsystem/internal/services"
	"time"
)

// Student Related Functions
// ------------------------------------------------------------------
func (r *RFIDRepository) GetStudentByRFID(studentId string) (*model.Student, error) {
	query := `
	SELECT student_ID, department_ID, first_Name, last_Name, middle_Name, birthday, contact_number, email, year_Level, program, block_section, first_access_timestamp, last_access_timestamp
	FROM Students
	WHERE student_ID = ?
	`

	log.Printf("Executing query with student ID: %s\n", studentId)

	student := &model.Student{}
	var firstAccessRaw, lastAccessRaw []byte
	err := r.dbClient.DB.QueryRow(query, studentId).Scan(
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
		&firstAccessRaw,
		&lastAccessRaw,
	)

	if err == sql.ErrNoRows {
		log.Printf("No student found with ID: %s\n", studentId)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error querying student: %v\n", err)
		return nil, fmt.Errorf("error querying student: %v", err)
	}

	// Parse raw timestamp bytes into *time.Time
	if len(firstAccessRaw) > 0 {
		if t, err := time.Parse("2006-01-02 15:04:05", string(firstAccessRaw)); err != nil {
			log.Printf("Invalid first_access_timestamp for %s: %v", student.StudentID, err)
		} else {
			student.FirstAccessTimestamp = &t
		}
	}
	if len(lastAccessRaw) > 0 {
		if t2, err := time.Parse("2006-01-02 15:04:05", string(lastAccessRaw)); err != nil {
			log.Printf("Invalid last_access_timestamp for %s: %v", student.StudentID, err)
		} else {
			student.LastAccessTimestamp = &t2
		}
	}

	// Update access timestamps
	now := time.Now()
	updateQuery := `
	UPDATE Students
	SET last_access_timestamp = ?,
		first_access_timestamp = COALESCE(first_access_timestamp, ?)
	WHERE student_ID = ?
	`
	if _, err := r.dbClient.DB.Exec(updateQuery, now, now, student.StudentID); err != nil {
		log.Printf("Error updating access timestamps for student %s: %v", student.StudentID, err)
	} else {
		if student.FirstAccessTimestamp == nil {
			student.FirstAccessTimestamp = &now
		}
		student.LastAccessTimestamp = &now
	}

	return student, nil
}

func (r *RFIDRepository) GetStudentSummaryData(studentId string) (*model.StudentInfoViewModel, error) {
	student, err := r.GetStudentByRFID(studentId)
	if err != nil {
		log.Printf("Error getting student summary data: %v\n", err)
		return nil, fmt.Errorf("error getting student summary data: %v", err)
	}

	if student == nil {
		log.Printf("No student found with ID: %s\n", studentId)
		return nil, nil
	}

	log.Printf("Student summary data retrieved for ID: %s\n", studentId)

	if student.YearLevel != nil {
		log.Printf("Student Year: %d", *student.YearLevel)
	} else {
		log.Printf("Student Year is nil")
	}

	var yearLevel string
	if student.YearLevel != nil {
		yearLevel = services.GetYearLevelString(*student.YearLevel)
	} else {
		yearLevel = ""
	}

	if student.YearLevel != nil {
		yearLevel = services.GetYearLevelString(*student.YearLevel)
	}

	gradesSummary, err := r.getStudentGradesSummary(studentId)
	if err != nil {
		log.Printf("Error getting grades summary: %v", err)
		// Continue with empty grades
	}

	// Get assessment
	assessment, err := r.getAssessment(studentId)
	if err != nil {
		log.Printf("Error getting assessment: %v", err)
		// Continue with empty assessment
	}

	// log.Printf("Student Assessment: %v", assessment)
	// Get payment schedules if we have an assessment
	var paymentSchedules []model.PaymentSchedule
	if assessment != nil {
		paymentSchedules, err = r.getPaymentSchedules(assessment.ID)
		log.Println("Pschedules", paymentSchedules)
		if err != nil {
			log.Printf("Error getting payment schedules: %v", err)
			// Continue with empty payment schedules
		}
	}

	log.Printf("Payment Schedules: %v", paymentSchedules)
	log.Printf("Payment Schedules Length: %d", len(paymentSchedules))
	log.Printf("Payment Schedules Type: %T", paymentSchedules)
	// log.Printf("Payment Schedules First Element Type: %T", paymentSchedules[0])

	return &model.StudentInfoViewModel{
		Student:          student,
		YearLevel:        yearLevel,
		GradesSummary:    gradesSummary,
		Assessment:       assessment,
		PaymentSchedules: paymentSchedules,
	}, nil

}

func (r *RFIDRepository) getStudentGradesSummary(studentId string) ([]model.YearGradeSummary, error) {
	log.Printf("Getting grades summary for student ID: %s", studentId)

	// Get current term to get the academic year
	currentTerm, err := r.GetCurrentTerm()
	if err != nil {
		log.Printf("Error getting current term: %v", err)
		return nil, err
	}
	if currentTerm == nil {
		log.Printf("No current term found")
		return nil, nil
	}

	log.Printf("Current academic year: %s", currentTerm.AcademicYear)

	// Add debug logs for query execution
	log.Printf("Executing grades summary query for academic year: %s", currentTerm.AcademicYear)
	query := `
	    SELECT
	    at.academic_year,
        at.semester,
        AVG(e.final_grade) as average_grade
    FROM
        Enrollments e
    JOIN
        AcademicTerms at ON e.term_id = at.term_id
    WHERE
        e.student_id = ?
        AND at.academic_year = ?
        AND e.final_grade IS NOT NULL
    GROUP BY
        at.academic_year, at.semester
    ORDER BY
        at.semester
    `

	rows, err := r.dbClient.DB.Query(query, studentId, currentTerm.AcademicYear)
	if err != nil {
		log.Printf("Error querying grades summary: %v", err)
		return nil, err
	}
	defer rows.Close()

	log.Printf("Query executed successfully, processing rows")

	// Create the year summary with empty grades
	emptyStr := "0.00"
	yearSummary := &model.YearGradeSummary{
		YearName:  currentTerm.AcademicYear,
		FirstSem:  &emptyStr,
		SecondSem: &emptyStr,
	}

	log.Printf("Created year summary with empty grades: %+v", yearSummary)

	for rows.Next() {
		var academicYear string
		var semester string
		var avgGrade float64

		if err := rows.Scan(&academicYear, &semester, &avgGrade); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		log.Printf("Scanned row - Year: %s, Semester: [%s], Grade: %f", academicYear, semester, avgGrade)

		// Format the grade as a string with 2 decimal places
		gradeStr := fmt.Sprintf("%.2f", avgGrade)
		// Create a new copy of the string to ensure it's not being shared/overwritten
		gradeCopy := string(gradeStr)
		log.Printf("Formatted grade: %s (copy: %s)", gradeStr, gradeCopy)

		// Log the semester value and its length for debugging
		log.Printf("Comparing semester value [%s] (length: %d)", semester, len(semester))

		// Set the semester grade based on exact string match
		switch semester {
		case "First Semester":
			log.Printf("Matched 'First' semester")
			var oldGrade string
			if yearSummary.FirstSem != nil {
				oldGrade = *yearSummary.FirstSem
			}
			log.Printf("Setting First semester grade from %v to %s", oldGrade, gradeCopy)
			yearSummary.FirstSem = &gradeCopy
		case "Second Semester":
			log.Printf("Matched 'Second' semester")
			var oldGrade string
			if yearSummary.SecondSem != nil {
				oldGrade = *yearSummary.SecondSem
			}
			log.Printf("Setting Second semester grade from %v to %s", oldGrade, gradeCopy)
			yearSummary.SecondSem = &gradeCopy
		default:
			log.Printf("Unknown semester value: %s", semester)
		}
	}
	// Final summary log with detailed grade values
	log.Printf("Final grades state - Academic Year: %s", yearSummary.YearName)
	if yearSummary.FirstSem != nil {
		log.Printf("First Semester: [%s]", *yearSummary.FirstSem)
	} else {
		log.Printf("First Semester: nil")
	}
	if yearSummary.SecondSem != nil {
		log.Printf("Second Semester: [%s]", *yearSummary.SecondSem)
	} else {
		log.Printf("Second Semester: nil")
	}

	// Create and return the result slice
	result := []model.YearGradeSummary{*yearSummary}
	log.Printf("Returning %d grade summaries", len(result))
	return result, nil
}
