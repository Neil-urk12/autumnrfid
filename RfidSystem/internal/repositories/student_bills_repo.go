package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"rfidsystem/internal/model"
	"time"
)

// Get payment schedules for an assessment
func (r *RFIDRepository) getPaymentSchedules(assessmentId int64) ([]model.PaymentSchedule, error) {
	query := "CALL getPaymentSchedules(?)"

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

// GetStudentBillsByRFID retrieves all billing-related data for a student by their RFID.
// This includes their assessment, fee breakdown, discounts, and payment history.
// It returns a Bills struct containing all this information or an error.
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
		log.Printf("No assessment found for student ID: %s", studentId)
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

	query := "CALL GetAssessment(?)"

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

	stmt, err := r.dbClient.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare fee breakdown query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(assessmentId)
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

	stmt, err := r.dbClient.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare discounts query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(assessmentId)
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
		payment_method,
		reference_number
	FROM Payments
	WHERE assessment_number = ?
	ORDER BY payment_date DESC
	`

	stmt, err := r.dbClient.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("prepare payment history query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(assessmentId)
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
			&payment.PaymentMethod,
			&payment.ReferenceNumber,
		); err != nil {
			return nil, err
		}
		// Format payment date to MM-DD-YYYY
		var t time.Time
		var err error
		t, err = time.Parse(time.RFC3339, payment.PaymentDate)
		if err != nil {
			// Fallback to date-only
			t, err = time.Parse("01-02-2006", payment.PaymentDate)
		}
		if err != nil {
			log.Printf("Invalid Payment Date format for payment: %v", err)
		} else {
			payment.PaymentDate = t.Format("01-02-2006")
		}

		payments = append(payments, payment)
	}

	return payments, nil
}
