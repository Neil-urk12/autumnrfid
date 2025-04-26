package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/model"
)

// Grades Related Functions
// ------------------------------------------------------------------

// GetStudentGradesByRFID retrieves grades for the current term
func (r *RFIDRepository) GetStudentGradesByRFID(studentId string) (*model.Grades, error) {
	currentTerm, err := r.GetCurrentTerm()
	if err != nil {
		return nil, fmt.Errorf("error getting current term: %v", err)
	}

	return r.GetStudentGradesByRFIDAndSemester(studentId, currentTerm.AcademicYear, currentTerm.Semester)
}

// GetStudentGradesByRFIDAndSemester retrieves grades for a specific semester in the given academic year
func (r *RFIDRepository) GetStudentGradesByRFIDAndSemester(studentId, academicYear, semesterName string) (*model.Grades, error) {
	student, err := r.GetStudentByRFID(studentId)
	if err != nil {
		log.Printf("Error getting student: %v\n", err)
		return nil, fmt.Errorf("error getting student: %v", err)
	}
	if student == nil {
		log.Printf("Student not found: %s\n", studentId)
		return nil, fmt.Errorf("student not found")
	}

	// Get term ID for the specified academic year and semester
	var termId int64
	termQuery := `
    SELECT term_id
    FROM AcademicTerms
    WHERE academic_year = ? AND semester = ?
    LIMIT 1`

	err = r.dbClient.DB.QueryRow(termQuery, academicYear, semesterName).Scan(&termId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no term found for academic year %s and semester %s", academicYear, semesterName)
	}
	if err != nil {
		return nil, fmt.Errorf("error finding term: %v", err)
	}

	query := `
    SELECT
        e.subject_Code,
        s.subject_name,
        e.prelim_grade,
        e.midterm_grade,
        e.prefinal_grade,
        e.final_term_grade,
        e.final_grade
    FROM Enrollments e
    JOIN Subjects s ON e.subject_Code = s.subject_Code
    JOIN AcademicTerms at ON e.term_id = at.term_id
    WHERE e.student_id = ?
        AND at.academic_year = ?
        AND at.semester = ?
    ORDER BY s.subject_code`

	rows, err := r.dbClient.DB.Query(query, studentId, academicYear, semesterName)
	if err != nil {
		return nil, fmt.Errorf("error querying grades: %v", err)
	}
	defer rows.Close()

	var gradeRecords []model.GradesRecord
	for rows.Next() {
		var grade model.GradesRecord
		err := rows.Scan(
			&grade.SubjectCode,
			&grade.SubjectName,
			&grade.PrelimGrade,
			&grade.MidtermGrade,
			&grade.PrefinalGrade,
			&grade.FinalTermGrade,
			&grade.FinalGrade,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning grade record: %v", err)
		}
		gradeRecords = append(gradeRecords, grade)
	}

	term := &model.AcademicTerm{
		ID:           termId,
		AcademicYear: academicYear,
		Semester:     semesterName,
	}

	return &model.Grades{
		Student:     student,
		CurrentTerm: term,
		Grades:      gradeRecords,
	}, nil
}

func (r *RFIDRepository) GetCurrentTerm() (*model.AcademicTerm, error) {
	query := `
	SELECT
		term_id,
		academic_year,
		semester,
		DATE_FORMAT(start_date, '%Y-%m-%d') as start_date,
        DATE_FORMAT(end_date, '%Y-%m-%d') as end_date
	FROM AcademicTerms
	WHERE CURRENT_DATE BETWEEN start_date AND end_date
	LIMIT 1
	`

	term := &model.AcademicTerm{}
	err := r.dbClient.DB.QueryRow(query).Scan(
		&term.ID,
		&term.AcademicYear,
		&term.Semester,
		&term.StartDate,
		&term.EndDate,
	)

	if err == sql.ErrNoRows {
		log.Printf("No current term found")
		return nil, nil
	}
	if err != nil {
		log.Printf("Error scanning current term: %v", err)
		return nil, err
	}

	log.Printf("Retrieved current term: ID=%d, Year=%s, Semester=%s",
		term.ID, term.AcademicYear, term.Semester)
	return term, nil
}
