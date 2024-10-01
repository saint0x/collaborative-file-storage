package scripts

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/saint0x/file-storage-app/backend/internal/db"
)

// SampleData holds all the sample data for populating the database
type SampleData struct {
	Users                    []User
	Friends                  []Friend
	FriendContexts           []FriendContext
	FriendLikes              []FriendLike
	Collections              []Collection
	Folders                  []Folder
	Files                    []File
	FileCategories           []FileCategory
	FileCategoryAssociations []FileCategoryAssociation
	SharedFiles              []SharedFile
	ActivityLog              []ActivityLog
}

// Struct definitions for all tables
type User struct {
	ID        string
	Email     string
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Friend struct {
	ID        string
	UserID    string
	FriendID  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FriendContext struct {
	ID        string
	UserID    string
	FriendID  string
	Context   string
	CreatedAt time.Time
}

type FriendLike struct {
	ID        string
	UserID    string
	FriendID  string
	CreatedAt time.Time
}

type Collection struct {
	ID          string
	UserID      string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Folder struct {
	ID          string
	UserID      string
	Name        string
	Description string
	ParentID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type File struct {
	ID           string
	UserID       string
	FolderID     string
	CollectionID string
	Key          string
	Name         string
	ContentType  string
	Size         int64
	UploadedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FileCategory struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FileCategoryAssociation struct {
	FileID     string
	CategoryID string
}

type SharedFile struct {
	ID         string
	FileID     string
	SharedBy   string
	SharedWith string
	CreatedAt  time.Time
}

type ActivityLog struct {
	ID            string
	UserID        string
	ActionType    string
	ActionDetails string
	CreatedAt     time.Time
}

// generateSampleData creates a set of sample data
func generateSampleData() SampleData {
	now := time.Now()
	return SampleData{
		Users: []User{
			{ID: uuid.New().String(), Email: "john@example.com", Username: "johndoe", FirstName: "John", LastName: "Doe", CreatedAt: now, UpdatedAt: now},
			{ID: uuid.New().String(), Email: "jane@example.com", Username: "janedoe", FirstName: "Jane", LastName: "Doe", CreatedAt: now, UpdatedAt: now},
		},
		Friends: []Friend{
			{ID: uuid.New().String(), UserID: "user1", FriendID: "user2", Status: "accepted", CreatedAt: now, UpdatedAt: now},
		},
		// Add more sample data for other tables as needed
	}
}

// PopulateSampleData is the main function to populate the database with sample data
func PopulateSampleData() {
	fmt.Println("🚀 Starting database population...")

	dbClient, err := db.NewSQLiteClient("./myapp.db")
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	sampleData := generateSampleData()

	if err := populateDatabase(dbClient, &sampleData); err != nil {
		log.Fatalf("❌ Failed to populate database: %v", err)
	}

	fmt.Println("✅ Database population complete!")
}

// populateDatabase inserts all sample data into the database
func populateDatabase(dbClient *db.SQLiteClient, data *SampleData) error {
	tx, err := dbClient.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := insertUsers(tx, data.Users); err != nil {
		return fmt.Errorf("failed to insert users: %w", err)
	}
	fmt.Println("👤 Users inserted")

	if err := insertFriends(tx, data.Friends); err != nil {
		return fmt.Errorf("failed to insert friends: %w", err)
	}
	fmt.Println("🤝 Friends inserted")

	// Add calls to other insert functions here

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Insert functions for each table
func insertUsers(tx *sql.Tx, users []User) error {
	stmt, err := tx.Prepare(`
		INSERT INTO users (id, email, username, first_name, last_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, user := range users {
		_, err := stmt.Exec(user.ID, user.Email, user.Username, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertFriends(tx *sql.Tx, friends []Friend) error {
	stmt, err := tx.Prepare(`
		INSERT INTO friends (id, user_id, friend_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, friend := range friends {
		_, err := stmt.Exec(friend.ID, friend.UserID, friend.FriendID, friend.Status, friend.CreatedAt, friend.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add more insert functions for other tables here...

// Helper function to generate a new UUID string
// func newUUID() string {
// 	return uuid.New().String()
// }
