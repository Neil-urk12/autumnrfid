package repositories

// RFIDRepository represents a repository for interacting with RFID-related data in the database.
// It encapsulates database operations related to students, bills, grades but
type RFIDRepository struct {
	dbClient *DatabaseClient
}

// NewRFIDRepository creates a new RFIDRepository with the provided DatabaseClient.
func NewRFIDRepository(dbClient *DatabaseClient) *RFIDRepository {
	return &RFIDRepository{dbClient: dbClient}
}

// Things to consider soon for readability and maintainability
// Break this down into sub repos
// student_repo.go
// student_bill_repo.go
// student_access_repo.go
// - to future me
