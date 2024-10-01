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

	_, err := c.DB.Exec(`
		INSERT INTO folders (id, user_id, name, description, parent_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, folder.ID, folder.UserID, folder.Name, folder.Description, folder.ParentID, folder.CreatedAt, folder.UpdatedAt)

	if err != nil {
		return "", err
	}

	return folder.ID.String(), nil
}

func (c *SQLiteClient) UpdateFileFolder(fileName string, folderID string) error {
	_, err := c.DB.Exec("UPDATE files SET folder_id = ? WHERE name = ?", folderID, fileName)
	return err
}

// Add these methods to the SQLiteClient struct

func (c *SQLiteClient) GetFileCategories() ([]string, error) {
	rows, err := c.DB.Query("SELECT name FROM file_categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *SQLiteClient) GetFilesByCategory(categoryName string) ([]models.File, error) {
	query := `
		SELECT f.id, f.user_id, f.name, f.content_type
		FROM files f
		JOIN file_category_associations fca ON f.id = fca.file_id
		JOIN file_categories fc ON fca.category_id = fc.id
		WHERE fc.name = ?
	`
	rows, err := c.DB.Query(query, categoryName)
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

func (c *SQLiteClient) GetFileDetails(fileID string) (models.FileDetails, error) {
	query := `
		SELECT f.id, f.user_id, f.name, f.content_type, f.key, f.size, f.uploaded_at, f.created_at, f.updated_at,
			   c.name as collection_name, fo.name as folder_name
		FROM files f
		LEFT JOIN collections c ON f.collection_id = c.id
		LEFT JOIN folders fo ON f.folder_id = fo.id
		WHERE f.id = ?
	`

	var details models.FileDetails
	var idStr, userIDStr string
	err := c.DB.QueryRow(query, fileID).Scan(
		&idStr, &userIDStr, &details.Name, &details.ContentType, &details.Key, &details.Size,
		&details.UploadedAt, &details.CreatedAt, &details.UpdatedAt,
		&details.CollectionName, &details.FolderName,
	)

	if err != nil {
		return models.FileDetails{}, err
	}

	details.ID, _ = uuid.Parse(idStr)
	details.UserID, _ = uuid.Parse(userIDStr)

	return details, nil
}

func (c *SQLiteClient) ShareFileWithFriends(fileID string, friendIDs []string) error {
	tx, err := c.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO shared_files (id, file_id, shared_by, shared_with) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, friendID := range friendIDs {
		_, err := stmt.Exec(uuid.New().String(), fileID, fileID, friendID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (c *SQLiteClient) GetSharedWithMeFiles(userID string) ([]models.File, error) {
	query := `
		SELECT f.id, f.user_id, f.name, f.content_type
		FROM files f
		JOIN shared_files sf ON f.id = sf.file_id
		WHERE sf.shared_with = ?
	`
	rows, err := c.DB.Query(query, userID)
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

func (c *SQLiteClient) GetOrganizedFileStructure(userID string) (models.FileStructure, error) {
	query := `
		WITH RECURSIVE folder_tree AS (
			SELECT id, name, parent_id, 0 AS level
			FROM folders
			WHERE user_id = ? AND parent_id IS NULL
			UNION ALL
			SELECT f.id, f.name, f.parent_id, ft.level + 1
			FROM folders f
			JOIN folder_tree ft ON f.parent_id = ft.id
		)
		SELECT ft.id, ft.name, ft.parent_id, ft.level, f.id as file_id, f.name as file_name
		FROM folder_tree ft
		LEFT JOIN files f ON f.folder_id = ft.id
		ORDER BY ft.level, ft.name, f.name
	`
	rows, err := c.DB.Query(query, userID)
	if err != nil {
		return models.FileStructure{}, err
	}
	defer rows.Close()

	structure := models.FileStructure{
		Folders: make(map[string]models.Folder),
		Files:   make(map[string][]models.File),
	}

	for rows.Next() {
		var folderID, folderName, parentID sql.NullString
		var level int
		var fileID, fileName sql.NullString

		err := rows.Scan(&folderID, &folderName, &parentID, &level, &fileID, &fileName)
		if err != nil {
			return models.FileStructure{}, err
		}

		if folderID.Valid {
			folder := models.Folder{
				ID:       uuid.MustParse(folderID.String),
				Name:     folderName.String,
				ParentID: uuid.NullUUID{UUID: uuid.MustParse(parentID.String), Valid: parentID.Valid},
			}
			structure.Folders[folderID.String] = folder
		}

		if fileID.Valid {
			file := models.File{
				ID:   uuid.MustParse(fileID.String),
				Name: fileName.String,
			}
			structure.Files[folderID.String] = append(structure.Files[folderID.String], file)
		}
	}

	return structure, nil
}

func (c *SQLiteClient) GetFriendContexts(friendID string) ([]string, error) {
	query := "SELECT context FROM friend_contexts WHERE friend_id = ?"
	rows, err := c.DB.Query(query, friendID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contexts []string
	for rows.Next() {
		var context string
		if err := rows.Scan(&context); err != nil {
			return nil, err
		}
		contexts = append(contexts, context)
	}
	return contexts, nil
}

func (c *SQLiteClient) AddFriendContext(userID, friendID, context string) error {
	_, err := c.DB.Exec("INSERT INTO friend_contexts (id, user_id, friend_id, context) VALUES (?, ?, ?, ?)",
		uuid.New().String(), userID, friendID, context)
	return err
}

func (c *SQLiteClient) RemoveFriendContext(userID, friendID, context string) error {
	_, err := c.DB.Exec("DELETE FROM friend_contexts WHERE user_id = ? AND friend_id = ? AND context = ?",
		userID, friendID, context)
	return err
}

func (c *SQLiteClient) LikeFriend(userID, friendID string) error {
	_, err := c.DB.Exec("INSERT INTO friend_likes (id, user_id, friend_id) VALUES (?, ?, ?)",
		uuid.New().String(), userID, friendID)
	return err
}

func (c *SQLiteClient) UnlikeFriend(userID, friendID string) error {
	_, err := c.DB.Exec("DELETE FROM friend_likes WHERE user_id = ? AND friend_id = ?", userID, friendID)
	return err
}

// Add this new struct to represent activity
type Activity struct {
	ID            string
	UserID        string
	ActionType    string
	ActionDetails string
	CreatedAt     time.Time
}

// Add this method to the SQLiteClient struct
func (c *SQLiteClient) GetRecentActivity(userID string) ([]Activity, error) {
	query := `
		SELECT id, user_id, action_type, action_details, created_at
		FROM activity_log
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := c.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent activity: %w", err)
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		err := rows.Scan(&activity.ID, &activity.UserID, &activity.ActionType, &activity.ActionDetails, &activity.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity row: %w", err)
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activity rows: %w", err)
	}

	return activities, nil
}

// Update the CreateFile function to include the b2_file_id
func (c *SQLiteClient) CreateFile(file models.File) error {
	_, err := c.DB.Exec(`
		INSERT INTO files (id, user_id, folder_id, collection_id, key, name, content_type, size, uploaded_at, created_at, updated_at, b2_file_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, file.ID, file.UserID, file.FolderID, file.CollectionID, file.Key, file.Name, file.ContentType, file.Size, file.UploadedAt, file.CreatedAt, file.UpdatedAt, file.B2FileID)
	return err
}

// Update the GetFileByID function to include the b2_file_id
func (c *SQLiteClient) GetFileByID(id string) (models.File, error) {
	var file models.File
	err := c.DB.QueryRow(`
		SELECT id, user_id, folder_id, collection_id, key, name, content_type, size, uploaded_at, created_at, updated_at, b2_file_id
		FROM files WHERE id = ?
	`, id).Scan(&file.ID, &file.UserID, &file.FolderID, &file.CollectionID, &file.Key, &file.Name, &file.ContentType, &file.Size, &file.UploadedAt, &file.CreatedAt, &file.UpdatedAt, &file.B2FileID)
	if err != nil {
		return models.File{}, err
	}
	return file, nil
}

// Update other relevant functions to include the b2_file_id field
