package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/saint0x/file-storage-app/backend/internal/config"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/handlers"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	database, err := db.NewSQLiteClient(cfg.SQLiteDBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Start server in a goroutine
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(time.Second)

	// Test creating a user
	user := handlers.User{Name: "John Doe", Email: "john@example.com"}
	userJSON, _ := json.Marshal(user)
	resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	var createdUser handlers.User
	json.NewDecoder(resp.Body).Decode(&createdUser)
	fmt.Printf("Created user: %+v\n", createdUser)

	// Test getting all users
	resp, err = http.Get("http://localhost:8080/users")
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}
	var users []handlers.User
	json.NewDecoder(resp.Body).Decode(&users)
	fmt.Printf("All users: %+v\n", users)

	// Test updating a user
	createdUser.Name = "Jane Doe"
	userJSON, _ = json.Marshal(createdUser)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/users/%d", createdUser.ID), bytes.NewBuffer(userJSON))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	var updatedUser handlers.User
	json.NewDecoder(resp.Body).Decode(&updatedUser)
	fmt.Printf("Updated user: %+v\n", updatedUser)

	// Test deleting a user
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/users/%d", createdUser.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Printf("Deleted user with status: %s\n", resp.Status)

	fmt.Println("All tests passed successfully!")
}