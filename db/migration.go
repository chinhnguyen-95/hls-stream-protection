package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

// MigrateJSONToDB migrates data from a JSON file to a database table
func MigrateJSONToDB(db *sql.DB, jsonFilePath, tableName string) error {
	// Read JSON file
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Parse JSON into a slice of maps
	var records []map[string]interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate records
	if len(records) == 0 {
		return fmt.Errorf("no records found in JSON file")
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Prepare insert statement
	columns := getColumns(records[0])
	query := buildInsertQuery(tableName, columns)
	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	// Insert records
	for _, record := range records {
		values := getValues(columns, record)
		if _, err := stmt.Exec(values...); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// getColumns extracts column names from a record
func getColumns(record map[string]interface{}) []string {
	columns := make([]string, 0, len(record))
	for column := range record {
		columns = append(columns, column)
	}
	return columns
}

// getValues extracts values from a record based on columns
func getValues(columns []string, record map[string]interface{}) []interface{} {
	values := make([]interface{}, len(columns))
	for i, column := range columns {
		values[i] = record[column]
	}
	return values
}

// buildInsertQuery creates an SQL insert query based on table name and columns
func buildInsertQuery(tableName string, columns []string) string {
	columnsList := ""
	placeholders := ""
	for i, column := range columns {
		if i > 0 {
			columnsList += ", "
			placeholders += ", "
		}
		columnsList += column
		placeholders += "?"
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnsList, placeholders)
}
