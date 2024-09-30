package main

import (
	"context"
	"log"
	"os"

	"github.com/saint0x/file-storage-app/internal/services/auth"
)

func main() {
	// Set the CLERK_SECRET_KEY environment variable
	os.Setenv("CLERK_SECRET_KEY", "your_clerk_secret_key_here")

	client, err := auth.NewClerkService()
	if err != nil {
		log.Fatalf("Failed to initialize Clerk client: %v", err)
	}

	// Example token verification
	token := "valid_token" // Replace with an actual token in a real scenario
	ctx := context.Background()
	userID, err := client.ValidateAndExtractUserID(ctx, token)
	if err != nil {
		log.Fatalf("Token verification failed: %v", err)
	}

	log.Printf("Token verified for user: %s", userID)
}
