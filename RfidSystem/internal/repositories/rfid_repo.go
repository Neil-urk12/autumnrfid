package repositories

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
