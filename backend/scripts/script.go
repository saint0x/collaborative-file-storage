package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
)

const (
	successEmoji = "✅"
	errorEmoji   = "❌"
	infoEmoji    = "ℹ️"
)

func main() {
	// Get the absolute path to the project root
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		logError("Failed to get project root path", err)
		os.Exit(1)
	}

	// Load environment variables
	if err := godotenv.Load(filepath.Join(projectRoot, ".env.local")); err != nil {
		logError("Error loading .env file", err)
		os.Exit(1)
	}
	logSuccess("Environment variables loaded")

	// Ensure server is running
	ensureServerIsRunning()

	// Initialize B2 service
	b2Service, err := storage.NewB2Service(
		os.Getenv("BACKBLAZE_KEY_ID"),
		os.Getenv("BACKBLAZE_APPLICATION_KEY"),
		os.Getenv("BACKBLAZE_BUCKET_NAME"),
		os.Getenv("BACKBLAZE_ENDPOINT"),
		os.Getenv("BACKBLAZE_REGION"),
	)
	if err != nil {
		logError("Failed to initialize B2 service", err)
		os.Exit(1)
	}
	logSuccess("B2 service initialized")

	// Upload test.json to B2 bucket
	if err := uploadTestJSON(b2Service); err != nil {
		logError("Failed to upload test.json", err)
		os.Exit(1)
	}

	logSuccess("Script completed successfully")
}

func ensureServerIsRunning() {
	logInfo("Checking if server is running...")
	_, err := http.Get("http://localhost:8080/health")
	if err != nil {
		logError("Server is not running", err)
		os.Exit(1)
	}
	logSuccess("Server is running")
}

func uploadTestJSON(b2Service *storage.B2Service) error {
	// Get the absolute path to the project root
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		return fmt.Errorf("failed to get project root path: %v", err)
	}

	logInfo("Reading test.json file...")
	jsonData, err := os.ReadFile(filepath.Join(projectRoot, "test.json"))
	if err != nil {
		return fmt.Errorf("failed to read test.json: %v", err)
	}
	logSuccess("test.json file read successfully")

	key := fmt.Sprintf("test_%s.json", time.Now().Format("20060102150405"))
	logInfo(fmt.Sprintf("Uploading file with key: %s", key))

	ctx := context.Background()
	err = b2Service.UploadFile(ctx, key, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	logSuccess(fmt.Sprintf("Successfully uploaded test.json as %s", key))

	logInfo("Verifying uploaded content...")
	downloadedContent, err := b2Service.DownloadFile(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to download uploaded file: %v", err)
	}
	defer downloadedContent.Close()

	content, err := io.ReadAll(downloadedContent)
	if err != nil {
		return fmt.Errorf("failed to read downloaded content: %v", err)
	}

	var downloadedJSON, originalJSON map[string]interface{}
	if err := json.Unmarshal(content, &downloadedJSON); err != nil {
		return fmt.Errorf("failed to unmarshal downloaded JSON: %v", err)
	}
	if err := json.Unmarshal(jsonData, &originalJSON); err != nil {
		return fmt.Errorf("failed to unmarshal original JSON: %v", err)
	}

	if !jsonEqual(downloadedJSON, originalJSON) {
		return fmt.Errorf("downloaded content does not match the original")
	}
	logSuccess("Verified uploaded content matches the original")

	return nil
}

func jsonEqual(a, b map[string]interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func logSuccess(message string) {
	log.Printf("%s %s\n", successEmoji, message)
}

func logError(message string, err error) {
	log.Printf("%s %s: %v\n", errorEmoji, message, err)
}

func logInfo(message string) {
	log.Printf("%s %s\n", infoEmoji, message)
}
