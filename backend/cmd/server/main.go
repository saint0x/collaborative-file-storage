package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/saint0x/file-storage-app/internal/services/auth"
)

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	client, err := auth.NewClerkService()
	if err != nil {
		log.Fatalf("Failed to initialize Clerk client: %v", err)
	}

	token := os.Getenv("TEST_TOKEN")
	ctx := context.Background()
	userID, err := client.ValidateAndExtractUserID(ctx, token)
	if err != nil {
		log.Fatalf("Token verification failed: %v", err)
	}

	log.Printf("Token verified for user: %s", userID)

	if err := startServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfig() error {
	if os.Getenv("CLERK_SECRET_KEY") == "" {
		return fmt.Errorf("CLERK_SECRET_KEY environment variable is not set")
	}
	return nil
}

func startServer() error {
	log.Println("Server starting...")
	// TODO: Implement your server startup logic here
	return nil
}
