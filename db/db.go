package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
)

// Connect initializes a database connection
func Connect(databaseURL string) (*sql.DB, error) {
	// Open connection
	db, err := sql.Open(determineDriver(databaseURL), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// determineDriver returns the driver based on the database URL
func determineDriver(databaseURL string) string {
	if len(databaseURL) >= 8 && databaseURL[:8] == "postgres" {
		return "postgres"
	}
	return "mysql"
}
