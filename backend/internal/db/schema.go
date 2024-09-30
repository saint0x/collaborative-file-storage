package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings
)

func InitSchema(db *sql.DB) error {
	schemaSQL, err := ioutil.ReadFile(filepath.Join("internal", "db", "schema.sql"))
	if err != nil {
		return fmt.Errorf("failed to read schema file: %v", err)
	}

	statements := strings.Split(string(schemaSQL), ";")
	for _, statement := range statements {
		trimmed := strings.TrimSpace(statement)
		if trimmed != "" {
			_, err := db.Exec(trimmed)
			if err != nil {
				return fmt.Errorf("failed to execute schema statement: %v", err)
			}
		}
	}

	return nil
}