package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saint0x/file-storage-app/backend/internal/models"
)

type SQLiteClient struct {
	*sql.DB
}

func NewSQLiteClient(dbPath string) (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return &SQLiteClient{DB: db}, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		)
	`)
	return err
}

func (c *SQLiteClient) GetFilesByIDs(fileIDs []string) ([]models.File, error) {
	query := `SELECT id, user_id, name, content_type FROM files WHERE id IN (?` + strings.Repeat(",?", len(fileIDs)-1) + `)`

	args := make([]interface{}, len(fileIDs))
	for i, id := range fileIDs {
		args[i] = id
	}

	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		var idStr, userIDStr string
		err := rows.Scan(&idStr, &userIDStr, &file.Name, &file.ContentType)
		if err != nil {
			return nil, err
		}
		file.ID, _ = uuid.Parse(idStr)
		file.UserID, _ = uuid.Parse(userIDStr)
		files = append(files, file)
	}

	return files, nil
}

func (c *SQLiteClient) CreateFolder(folder models.Folder) (string, error) {
	folder.ID = uuid.New()
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	_, err := c.DB.Exec("INSERT INTO folders (id, user_id, name, description, parent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		folder.ID, folder.UserID, folder.Name, folder.Description, folder.ParentID, folder.CreatedAt, folder.UpdatedAt)
	if err != nil {
		return "", err
	}

	return folder.ID.String(), nil
}

func (c *SQLiteClient) UpdateFileFolder(fileName string, folderID string) error {
	_, err := c.DB.Exec("UPDATE files SET folder_id = ? WHERE name = ?", folderID, fileName)
	return err
}
