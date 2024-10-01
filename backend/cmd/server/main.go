package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
)

func main() {
	if err := loadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()

	// Initialize auth service
	authService, err := auth.NewClerkService()
	if err != nil {
		log.Fatalf("Failed to initialize Clerk client: %v", err)
	}

	// Validate token (for demonstration purposes)
	token := os.Getenv("TEST_TOKEN")
	userID, err := authService.ValidateAndExtractUserID(ctx, token)
	if err != nil {
		log.Fatalf("Token verification failed: %v", err)
	}
	log.Printf("Token verified for user: %s", userID)

	// Initialize B2 service
	b2Service, err := storage.NewB2Service(
		os.Getenv("BACKBLAZE_KEY_ID"),
		os.Getenv("BACKBLAZE_APPLICATION_KEY"),
		os.Getenv("BACKBLAZE_BUCKET_NAME"),
		os.Getenv("BACKBLAZE_ENDPOINT"),
		"us-east-001", // Hardcoded region
	)
	if err != nil {
		log.Fatalf("Failed to initialize B2 service: %v", err)
	}

	// Demonstrate B2Service functionality
	if err := demonstrateB2Service(ctx, b2Service); err != nil {
		log.Fatalf("B2 service demonstration failed: %v", err)
	}

	if err := startServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfig() error {
	requiredEnvVars := []string{
		"CLERK_SECRET_KEY",
		"BACKBLAZE_KEY_ID",
		"BACKBLAZE_APPLICATION_KEY",
		"BACKBLAZE_BUCKET_NAME",
		"BACKBLAZE_ENDPOINT",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%s environment variable is not set", envVar)
		}
	}
	return nil
}

func startServer() error {
	log.Println("Server starting...")
	// TODO: Implement your server startup logic here
	return nil
}

func demonstrateB2Service(ctx context.Context, b2Service *storage.B2Service) error {
	// Create a new bucket
	newBucketName := "test-bucket-" + time.Now().Format("20060102150405")
	if err := b2Service.CreateBucket(ctx, newBucketName); err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}
	log.Printf("Created new bucket: %s", newBucketName)

	// Set bucket ACL to public-read
	if err := b2Service.SetBucketACL(ctx, newBucketName, "public-read"); err != nil {
		return fmt.Errorf("failed to set bucket ACL: %v", err)
	}
	log.Printf("Set bucket ACL to public-read: %s", newBucketName)

	// Upload a test file
	testContent := strings.NewReader("This is a test file content")
	testKey := "test-file.txt"
	if err := b2Service.UploadFile(ctx, testKey, testContent); err != nil {
		return fmt.Errorf("failed to upload test file: %v", err)
	}
	log.Printf("Test file uploaded successfully: %s", testKey)

	// List files
	files, err := b2Service.ListFiles(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to list files: %v", err)
	}
	log.Printf("Files in bucket: %v", files)

	// Generate a signed URL
	signedURL, err := b2Service.GetSignedURL(ctx, testKey, 1*time.Hour)
	if err != nil {
		return fmt.Errorf("failed to generate signed URL: %v", err)
	}
	log.Printf("Signed URL for test file: %s", signedURL)

	// Download the test file
	downloadedContent, err := b2Service.DownloadFile(ctx, testKey)
	if err != nil {
		return fmt.Errorf("failed to download test file: %v", err)
	}
	defer downloadedContent.Close()
	content, err := io.ReadAll(downloadedContent)
	if err != nil {
		return fmt.Errorf("failed to read downloaded content: %v", err)
	}
	log.Printf("Downloaded file content: %s", string(content))

	// Delete the test file
	if err := b2Service.DeleteFile(ctx, testKey); err != nil {
		return fmt.Errorf("failed to delete test file: %v", err)
	}
	log.Printf("Test file deleted successfully: %s", testKey)

	// Delete the test bucket
	if err := b2Service.DeleteBucket(ctx, newBucketName); err != nil {
		return fmt.Errorf("failed to delete test bucket: %v", err)
	}
	log.Printf("Test bucket deleted successfully: %s", newBucketName)

	return nil
}
