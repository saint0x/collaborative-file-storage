package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

func InitSchema(db *sql.DB) error {
	// Check if the users table exists
	var tableName string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			// Users table doesn't exist, create all tables
			schemaSQL, err := ioutil.ReadFile(filepath.Join("internal", "db", "schema.sql"))
			if err != nil {
				return fmt.Errorf("failed to read schema file: %w", err)
			}

			_, err = db.Exec(string(schemaSQL))
			if err != nil {
				return fmt.Errorf("failed to execute schema SQL: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check if users table exists: %w", err)
		}
	}

	// Run migrations
	err = runMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations(db *sql.DB) error {
	files, err := filepath.Glob("internal/db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to list migration files: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Executing migration: %s\n", file)
		migrationSQL, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Split the file content into "Up" and "Down" migrations
		parts := strings.Split(string(migrationSQL), "-- Down migration")
		upMigration := strings.TrimSpace(parts[0])

		// Start a transaction for each migration
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
		}

		// Execute each statement in the migration separately
		statements := strings.Split(upMigration, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			_, err = tx.Exec(stmt)
			if err != nil {
				tx.Rollback()
				fmt.Printf("Error executing statement: %s\n", stmt)
				return fmt.Errorf("failed to execute statement in %s: %w\nStatement: %s", file, err, stmt)
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for %s: %w", file, err)
		}

		fmt.Printf("Successfully executed migration: %s\n", file)
	}

	return nil
}
