package model

import (
	"encoding/json"
	"strings"
	"time"
)

type StudentInfo struct {
	Student *Student
}

type YearGradeSummary struct {
	YearName  string  `json:"year_name"`
	FirstSem  *string `json:"first_sem,omitempty"`
	SecondSem *string `json:"second_sem,omitempty"`
}

type StudentInfoViewModel struct {
	Student          *Student
	YearLevel        string
	GradesSummary    []YearGradeSummary
	Assessment       *Assessment
	PaymentSchedules []PaymentSchedule
}

type Department struct {
	ID   int64  `json:"department_id" db:"department_ID"`
	Name string `json:"department_name" db:"department_Name"`
}

type Subject struct {
	SubjectCode string `json:"subject_code" db:"subject_code"`
	SubjectName string `json:"subject_name" db:"subject_name"`
}

type AcademicTerm struct {
	ID           int64  `json:"term_id" db:"term_id"`
	AcademicYear string `json:"academic_year" db:"academic_year"`
	// Semester     *time.Time `json:"semester" db:"semester"`
	Semester string `json:"semester" db:"semester"`
	// StartDate    *time.Time `json:"start_date,omitempty" db:"start_date"`
	// EndDate      *time.Time `json:"end_date,omitempty" db:"end_date"`
	StartDate string `json:"start_date,omitempty" db:"start_date"`
	EndDate   string `json:"end_date,omitempty" db:"end_date"`
}

type FeeType struct {
	ID       int64  `json:"fee_type_id" db:"fee_type_id"`
	Name     string `json:"name" db:"name"`
	Category string `json:"category" db:"category"`
}

type DiscountType struct {
	ID           int64   `json:"discount_type_id" db:"discount_type_id"`
	Name         string  `json:"name" db:"name"`
	Description  *string `json:"description,omitempty" db:"description"`
	IsPercentage bool    `json:"is_percentage" db:"is_percentage"`
	Value        float64 `json:"value" db:"value"`
}

type Student struct {
	StudentID            string     `json:"student_id" db:"student_ID"`
	DepartmentID         *int64     `json:"department_id,omitempty" db:"department_ID"`
	FirstName            *string    `json:"first_name,omitempty" db:"first_Name"`
	LastName             *string    `json:"last_name,omitempty" db:"last_Name"`
	MiddleName           *string    `json:"middle_name,omitempty" db:"middle_Name"`
	Birthday             *string    `json:"birthday,omitempty" db:"birthday"`
	ContactNumber        *string    `json:"contact_number,omitempty" db:"contact_number"`
	Email                *string    `json:"email,omitempty" db:"email"`
	YearLevel            *int       `json:"year_level,omitempty" db:"year_Level"`
	Program              *string    `json:"program,omitempty" db:"program"`
	BlockSection         *string    `json:"block_section,omitempty" db:"block_section"`
	FirstAccessTimestamp *time.Time `json:"first_access_timestamp,omitempty" db:"first_access_timestamp"`
	LastAccessTimestamp  *time.Time `json:"last_access_timestamp,omitempty" db:"last_access_timestamp"`
	// FirstAccessTimestamp string  `json:"first_access_timestamp" db:"first_access_timestamp"`
	// LastAccessTimestamp  string  `json:"last_access_timestamp" db:"last_access_timestamp"`
	Status *string `json:"status,omitempty" db:"status"`
}

// MarshalJSON customizes Student JSON to format timestamps as "YYYY-MM-DD hh:mm am/pm" without seconds
func (s *Student) MarshalJSON() ([]byte, error) {
	type Alias Student
	aux := &struct {
		*Alias
		FirstAccessTimestamp string `json:"first_access_timestamp,omitempty"`
		LastAccessTimestamp  string `json:"last_access_timestamp,omitempty"`
	}{
		Alias: (*Alias)(s),
	}
	if s.FirstAccessTimestamp != nil {
		aux.FirstAccessTimestamp = strings.ToLower(s.FirstAccessTimestamp.Format("2006-01-02 03:04 PM"))
	}
	if s.LastAccessTimestamp != nil {
		aux.LastAccessTimestamp = strings.ToLower(s.LastAccessTimestamp.Format("2006-01-02 03:04 PM"))
	}
	return json.Marshal(aux)
}

type Assessment struct {
	ID                  int64    `json:"assessment_number" db:"assessment_Number"`
	StudentID           *string  `json:"student_id,omitempty" db:"student_ID"`
	TermID              *int64   `json:"term_id,omitempty" db:"term_id"`
	TotalFeeAmount      float64  `json:"total_fee_amount" db:"total_fee_amount"`
	TotalDiscountAmount float64  `json:"total_discount_amount" db:"total_discount_amount"`
	NetAssessmentAmount float64  `json:"net_assessment_amount" db:"net_assessment_amount"`
	InitialPayment      *float64 `json:"initial_payment,omitempty" db:"initial_Payment"`
	TotalPaymentAmount  float64  `json:"total_payment_amount" db:"total_payment_amount"`
	FullPmtIfB4Prelim   *float64 `json:"full_pmt_if_b4_prelim,omitempty" db:"full_pmt_if_b4_prelim"`
	RemainingBalance    float64  `json:"remaining_balance" db:"remaining_Balance"`
	PerExamFee          *float64 `json:"per_exam_fee,omitempty" db:"per_Exam_Fee"`
}

type AssessmentViewModel struct {
	ID                  int64
	StudentID           *string
	TermID              *int64
	TotalFeeAmount      string
	NetAssessmentAmount string
	InitialPayment      string
	TotalPaymentAmount  string
	RemainingBalance    string
	TotalDiscountAmount string
	FullPmtIfB4Prelim   string
	PerExamFee          string
}

// TotalAmount == TotalAssessment
// NetAmount == Total Amount
// InitialPayment == InitialPayment
// RemainingBalance == Remaining Balance
// Old Continuing Discount == Total Discount

type Enrollment struct {
	ID          int64    `json:"enrollment_id" db:"enrollment_ID"`
	StudentID   *string  `json:"student_id,omitempty" db:"student_ID"`
	SubjectCode *string  `json:"subject_code,omitempty" db:"subject_Code"`
	TermID      *int64   `json:"term_id,omitempty" db:"term_id"`
	PrelimGrade *float64 `json:"prelim_grade,omitempty" db:"prelim_grade"`

	MidtermGrade   *float64 `json:"midterm_grade,omitempty" db:"midterm_grade"`
	PrefinalGrade  *float64 `json:"prefinal_grade,omitempty" db:"prefinal_grade"`
	FinalTermGrade *float64 `json:"final_term_grade,omitempty" db:"final_term_grade"`
	FinalGrade     *float64 `json:"final_grade,omitempty" db:"final_grade"`
}

type AssessmentFee struct {
	ID               int64   `json:"assessment_fee_id" db:"assessment_fee_id"`
	AssessmentNumber int64   `json:"assessment_number" db:"assessment_number"`
	FeeTypeID        int64   `json:"fee_type_id" db:"fee_type_id"`
	Amount           float64 `json:"amount" db:"amount"`
}

type AssessmentDiscount struct {
	ID               int64    `json:"assessment_discount_id" db:"assessment_discount_id"`
	AssessmentNumber int64    `json:"assessment_number" db:"assessment_number"`
	DiscountTypeID   int64    `json:"discount_type_id" db:"discount_type_id"`
	AppliedAmount    float64  `json:"applied_amount" db:"applied_amount"`
	CalculationBasis *float64 `json:"calculation_basis,omitempty" db:"calculation_basis"`
}

type Payment struct {
	ID               int64 `json:"payment_id" db:"payment_id"`
	AssessmentNumber int64 `json:"assessment_number" db:"assessment_number"`
	// PaymentDate      time.Time `json:"payment_date" db:"payment_date"`
	PaymentDate string  `json:"payment_data" db:"payment_data"`
	Amount      float64 `json:"amount" db:"amount"`
	Description *string `json:"description,omitempty" db:"description"`
	// Status          string  `json:"status" db:"status"`
	PaymentMethod   *string `json:"payment_method,omitempty" db:"payment_method"`
	ReferenceNumber *string `json:"reference_number,omitempty" db:"reference_number"`
}

type PaymentSchedule struct {
	ID               int64  `json:"schedule_id" db:"schedule_id"`
	AssessmentNumber int64  `json:"assessment_number" db:"assessment_number"`
	TermDescription  string `json:"term_description" db:"term_description"`
	// DueDate          *time.Time `json:"due_date,omitempty" db:"due_date"`
	DueDate        string  `json:"due_date" db:"due_date"`
	ExpectedAmount float64 `json:"expected_amount" db:"expected_amount"`
	SortOrder      int     `json:"sort_order" db:"sort_order"`
}

type PaymentScheduleViewModel struct {
	ID                      int64
	AssessmentNumber        int64
	TermDescription         string
	DueDate                 string
	ExpectedAmount          float64
	ExpectedAmountFormatted string
	SortOrder               int
}

// -------------------------
// Rfid Repo structs
type FeeBreakdown struct {
	Category string
	Name     string
	Amount   float64
}

type PaymentRecord struct {
	// Date            time.Time
	PaymentDate string
	Description *string
	Amount      float64
	// Status          string
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
	Assessment     *Assessment
	FeeBreakdown   []FeeBreakdown
	Discounts      []DiscountRecord
	PaymentHistory []PaymentRecord
}

type GradesRecord struct {
	SubjectCode    string   `json:"subject_code"`
	SubjectName    string   `json:"subject_name"`
	PrelimGrade    *float64 `json:"prelim_grade"`
	MidtermGrade   *float64 `json:"midterm_grade"`
	PrefinalGrade  *float64 `json:"prefinal_grade"`
	FinalGrade     *float64 `json:"final_grade"`
	FinalTermGrade *float64 `json:"final_term_grade"`
}

type Grades struct {
	Student     *Student
	CurrentTerm *AcademicTerm
	Grades      []GradesRecord
}

// StudentAssessmentSummary represents the summary of students for a specific assessment term.
type StudentAssessmentSummary struct {
	StudentID string  `json:"student_id" db:"StudentID"`
	Name      *string `json:"name,omitempty" db:"Name"`
	Course    *string `json:"course,omitempty" db:"Course"`
	YearLevel *string `json:"year_level,omitempty" db:"YearLevel"`
	Status    *string `json:"status,omitempty" db:"status"`
}

// PaginationMetadata holds information about the pagination state.
type PaginationMetadata struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
}

// PaginatedStudentAssessmentResponse structures the paginated response for student assessments.
type PaginatedStudentAssessmentResponse struct {
	Data       []*StudentAssessmentSummary `json:"data"`
	Pagination PaginationMetadata          `json:"pagination"`
}
