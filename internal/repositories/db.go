package repositories

import (
	"database/sql"
	"fmt"
	"rfidsystem/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseClient struct {
	DB *sql.DB
}

func NewDatabaseClient(config config.DatabaseConfig) (*DatabaseClient, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		config.Username, config.Password, config.Host, config.Port, config.DatabaseName,
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return &DatabaseClient{DB: db}, nil
}
