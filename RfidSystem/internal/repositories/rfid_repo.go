package repositories

import (
	"database/sql"
	"fmt"
	"rfidsystem/internal/model"
)

type FeeBreakdown struct {
	Category string
	Name     string
	Amount   float64
}

type PaymentRecord struct {
	// Date            time.Time
	PaymentDate     string
	Description     *string
	Amount          float64
	Status          string
	PaymentMethod   *string
	ReferenceNumber *string
}

type DiscountRecord struct {
	Name             string
	IsPercentage     bool
	Value            float64
	AppliedAmount    float64
	CalculationBasis *float64
}

type Bills struct {
	Assessment     *model.Assessment
	FeeBreakdown   []FeeBreakdown
	Discounts      []DiscountRecord
	PaymentHistory []PaymentRecord
}

type Grades struct {
	// Add grades data here
}

type RFIDRepository struct {
	dbClient *DatabaseClient
}

func NewRFIDRepository(dbClient *DatabaseClient) *RFIDRepository {
	return &RFIDRepository{dbClient: dbClient}
}

func (r *RFIDRepository) GetStudentByRFID(studentId string) (*model.Student, error) {
	query := `
	SELECT student_ID, department_ID, first_Name, last_Name, middle_Name, birthday, contact_number, email, year_Level, program, block_section, first_access_timestamp, last_access_timestamp
	FROM Students 
	WHERE student_ID = ?
	`

	fmt.Printf("Executing query with student ID: %s\n", studentId)

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
		fmt.Printf("No student found with ID: %s\n", studentId)
		return nil, nil
	}
	if err != nil {
		fmt.Printf("Error querying student: %v\n", err)
		return nil, fmt.Errorf("error querying student: %v", err)
	}

	return student, nil
}

func (r *RFIDRepository) GetStudentBillsByRFID(studentId string) (*Bills, error) {
	// Test database connection
	if err := r.dbClient.DB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection error: %v", err)
	}
	fmt.Printf("Database connection confirmed for student ID: %s\n", studentId)

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

	return &Bills{
		Assessment:     assessment,
		FeeBreakdown:   fees,
		Discounts:      discounts,
		PaymentHistory: payments,
	}, nil
}

func (r *RFIDRepository) getAssessment(studentId string) (*model.Assessment, error) {
	fmt.Printf("Getting assessment for student ID: %s\n", studentId)
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
	fmt.Printf("Executing assessment query with parameters: studentId=%s\n", studentId)
	row := r.dbClient.DB.QueryRow(query, studentId)
	if row == nil {
		return nil, fmt.Errorf("database returned nil row")
	}

	fmt.Printf("Scanning assessment row into struct...\n")
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
		fmt.Printf("No assessment found for student ID: %s\n", studentId)
		return nil, nil
	}
	if err != nil {
		fmt.Printf("Database error getting assessment: %v\n", err)
		return nil, fmt.Errorf("database error: %v", err)
	}
	fmt.Printf("Successfully retrieved assessment for student ID: %s\n", studentId)

	return assessment, nil
}

func (r *RFIDRepository) getFeeBreakdown(assessmentId int64) ([]FeeBreakdown, error) {
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

	var fees []FeeBreakdown
	for rows.Next() {
		var fee FeeBreakdown
		if err := rows.Scan(&fee.Category, &fee.Name, &fee.Amount); err != nil {
			return nil, err
		}
		fees = append(fees, fee)
	}

	return fees, nil
}

func (r *RFIDRepository) getDiscounts(assessmentId int64) ([]DiscountRecord, error) {
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

	var discounts []DiscountRecord
	for rows.Next() {
		var discount DiscountRecord
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

func (r *RFIDRepository) getPaymentHistory(assessmentId int64) ([]PaymentRecord, error) {
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

	var payments []PaymentRecord
	for rows.Next() {
		var payment PaymentRecord
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

func (r *RFIDRepository) GetStudentGradesByRFID(studentId string) (*Grades, error) {
	return nil, nil
}
