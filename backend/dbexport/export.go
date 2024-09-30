package dbexport

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saint0x/file-storage-app/backend/internal/config"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"gopkg.in/yaml.v2"
)

type TableData struct {
	Columns []string        `yaml:"columns"`
	Rows    [][]interface{} `yaml:"rows"`
}

type DatabaseContent struct {
	Tables map[string]TableData `yaml:"tables"`
}

func ExportDatabase() {
	log.Println("🚀 Starting database export script...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Connect to the database
	database, err := db.NewSQLiteClient(cfg.SQLiteDBPath)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Get all table names
	tables, err := getTableNames(database.DB)
	if err != nil {
		log.Fatalf("❌ Failed to get table names: %v", err)
	}

	// Initialize DatabaseContent
	dbContent := DatabaseContent{
		Tables: make(map[string]TableData),
	}

	// For each table, get its structure and data
	for _, table := range tables {
		tableData, err := getTableData(database.DB, table)
		if err != nil {
			log.Fatalf("❌ Failed to get data for table %s: %v", table, err)
		}
		dbContent.Tables[table] = tableData
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(dbContent)
	if err != nil {
		log.Fatalf("❌ Failed to convert data to YAML: %v", err)
	}

	// Write to file
	err = os.WriteFile("your-db.yaml", yamlData, 0644)
	if err != nil {
		log.Fatalf("❌ Failed to write YAML file: %v", err)
	}

	log.Println("✅ Database content exported to your-db.yaml successfully!")
}

func getTableNames(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func getTableData(db *sql.DB, tableName string) (TableData, error) {
	// Get column names
	columns, err := getColumnNames(db, tableName)
	if err != nil {
		return TableData{}, err
	}

	// Get row data
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return TableData{}, err
	}
	defer rows.Close()

	var tableData TableData
	tableData.Columns = columns

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return TableData{}, err
		}

		tableData.Rows = append(tableData.Rows, values)
	}

	return tableData, nil
}

func getColumnNames(db *sql.DB, tableName string) ([]string, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var (
			cid        int
			name       string
			type_      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, nil
}
