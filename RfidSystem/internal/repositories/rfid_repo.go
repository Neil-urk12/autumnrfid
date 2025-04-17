package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"rfidsystem/internal/services"
)

type RFIDRepository struct {
	dbClient *DatabaseClient
}

func NewRFIDRepository(dbClient *DatabaseClient) *RFIDRepository {
	return &RFIDRepository{dbClient: dbClient}
}

// Things to consider soon for readability and maintainability
// Break this down into sub repos
// student_repo.go
// student_bill_repo.go
// student_access_repo.go
// - to future me

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
		&student.FirstAccessTimestamp,
		&student.LastAccessTimestamp,
	)

	if err == sql.ErrNoRows {
		log.Printf("No student found with ID: %s\n", studentId)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error querying student: %v\n", err)
		return nil, fmt.Errorf("error querying student: %v", err)
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
		return nil, fmt.Errorf("no student found with ID: %s", studentId)
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

// Get payment schedules for an assessment
func (r *RFIDRepository) getPaymentSchedules(assessmentId int64) ([]model.PaymentSchedule, error) {
	query := `
    SELECT
        schedule_id,
        assessment_number,
        term_description,
        due_date,
        expected_amount,
        sort_order
    FROM
        PaymentSchedule
    WHERE
        assessment_number = ?
    ORDER BY
        sort_order
    `

	rows, err := r.dbClient.DB.Query(query, assessmentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.PaymentSchedule
	for rows.Next() {
		var schedule model.PaymentSchedule
		if err := rows.Scan(
			&schedule.ID,
			&schedule.AssessmentNumber,
			&schedule.TermDescription,
			&schedule.DueDate,
			&schedule.ExpectedAmount,
			&schedule.SortOrder,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if schedules == nil {
		schedules = []model.PaymentSchedule{}
	}

	return schedules, nil
}

// Bills Related Functions
// ------------------------------------------------------------------

func (r *RFIDRepository) GetStudentBillsByRFID(studentId string) (*model.Bills, error) {
	// Test database connection
	if err := r.dbClient.DB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}
	log.Printf("Database connection confirmed for student ID: %s\n", studentId)

	// Get main assessment data
	assessment, err := r.getAssessment(studentId)
	if err != nil {
		return nil, fmt.Errorf("error getting assessment: %v", err)
	}
	if assessment == nil {
		return nil, nil
	}

	// Get fee breakdown
	fees, err := r.getFeeBreakdown(assessment.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting fee breakdown: %v", err)
	}

	// Get discounts
	discounts, err := r.getDiscounts(assessment.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting discounts: %v", err)
	}

	// Get payment history
	payments, err := r.getPaymentHistory(assessment.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting payment history: %v", err)
	}

	log.Printf("Payment history retrieved for student ID: %s", studentId)

	return &model.Bills{
		Assessment:     assessment,
		FeeBreakdown:   fees,
		Discounts:      discounts,
		PaymentHistory: payments,
	}, nil
}

func (r *RFIDRepository) getAssessment(studentId string) (*model.Assessment, error) {
	log.Printf("Getting assessment for student ID: %s\n", studentId)
	query := `
	SELECT
		assessment_Number,
		student_ID,
		term_id,
		total_fee_amount,
		total_discount_amount,
		net_assessment_amount,
		initial_Payment,
		total_payment_amount,
		full_pmt_if_b4_prelim,
		remaining_Balance,
		per_Exam_Fee
	FROM Assessment
	WHERE student_ID = ?
	ORDER BY assessment_Number DESC
	LIMIT 1
	`

	assessment := &model.Assessment{}
	log.Printf("Executing assessment query with parameters: studentId=%s\n", studentId)
	row := r.dbClient.DB.QueryRow(query, studentId)
	if row == nil {
		return nil, fmt.Errorf("database returned nil row")
	}

	log.Printf("Scanning assessment row into struct...\n")
	err := row.Scan(
		&assessment.ID,
		&assessment.StudentID,
		&assessment.TermID,
		&assessment.TotalFeeAmount,
		&assessment.TotalDiscountAmount,
		&assessment.NetAssessmentAmount,
		&assessment.InitialPayment,
		&assessment.TotalPaymentAmount,
		&assessment.FullPmtIfB4Prelim,
		&assessment.RemainingBalance,
		&assessment.PerExamFee,
	)

	if err == sql.ErrNoRows {
		log.Printf("No assessment found for student ID: %s\n", studentId)
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error getting assessment: %v\n", err)
		return nil, fmt.Errorf("database error: %v", err)
	}
	log.Printf("Successfully retrieved assessment for student ID: %s\n", studentId)

	return assessment, nil
}

func (r *RFIDRepository) getFeeBreakdown(assessmentId int64) ([]model.FeeBreakdown, error) {
	query := `
	SELECT
		ft.category,
		ft.name,
		af.amount
	FROM AssessmentFees af
	JOIN FeeTypes ft ON af.fee_type_id = ft.fee_type_id
	WHERE af.assessment_number = ?
	ORDER BY ft.category, ft.name
	`

	rows, err := r.dbClient.DB.Query(query, assessmentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fees []model.FeeBreakdown
	for rows.Next() {
		var fee model.FeeBreakdown
		if err := rows.Scan(&fee.Category, &fee.Name, &fee.Amount); err != nil {
			return nil, err
		}
		fees = append(fees, fee)
	}

	return fees, nil
}

func (r *RFIDRepository) getDiscounts(assessmentId int64) ([]model.DiscountRecord, error) {
	query := `
	SELECT
		dt.name,
		dt.is_percentage,
		dt.value,
		ad.applied_amount,
		ad.calculation_basis
	FROM AssessmentDiscounts ad
	JOIN DiscountTypes dt ON ad.discount_type_id = dt.discount_type_id
	WHERE ad.assessment_number = ?
	`

	rows, err := r.dbClient.DB.Query(query, assessmentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts []model.DiscountRecord
	for rows.Next() {
		var discount model.DiscountRecord
		if err := rows.Scan(
			&discount.Name,
			&discount.IsPercentage,
			&discount.Value,
			&discount.AppliedAmount,
			&discount.CalculationBasis,
		); err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}

	return discounts, nil
}

func (r *RFIDRepository) getPaymentHistory(assessmentId int64) ([]model.PaymentRecord, error) {
	log.Printf("Fetching payment history for assessment ID: %d", assessmentId)
	query := `
	SELECT
		payment_date,
		description,
		amount,
		status,
		payment_method,
		reference_number
	FROM Payments
	WHERE assessment_number = ?
	ORDER BY payment_date DESC
	`

	rows, err := r.dbClient.DB.Query(query, assessmentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.PaymentRecord
	for rows.Next() {
		var payment model.PaymentRecord
		if err := rows.Scan(
			&payment.PaymentDate,
			&payment.Description,
			&payment.Amount,
			&payment.Status,
			&payment.PaymentMethod,
			&payment.ReferenceNumber,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

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
