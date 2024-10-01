package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	"github.com/saint0x/file-storage-app/backend/internal/api"
	apimiddleware "github.com/saint0x/file-storage-app/backend/internal/api/middleware"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/services/ai"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
	"github.com/saint0x/file-storage-app/backend/pkg/logger"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Load configuration
	if err := loadConfig(); err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize database
	dbClient, err := db.NewSQLiteClient("file_storage.db")
	if err != nil {
		logger.Fatal("Failed to initialize database", "error", err)
	}
	defer dbClient.Close()

	// Initialize services
	authService, err := auth.NewClerkService()
	if err != nil {
		logger.Fatal("Failed to initialize Clerk service", "error", err)
	}

	b2Service, err := storage.NewB2Service(
		os.Getenv("BACKBLAZE_ACCOUNT_ID"),
		os.Getenv("BACKBLAZE_APPLICATION_KEY"),
		os.Getenv("BACKBLAZE_BUCKET_ID"),
	)
	if err != nil {
		logger.Fatal("Failed to initialize B2 service", "error", err)
	}

	wsHub := websocket.NewHub(dbClient)

	aiProcessor := ai.NewProcessor(os.Getenv("OPENAI_API_KEY"))

	// Set up router
	r := chi.NewRouter()

	// Middleware
	r.Use(apimiddleware.RequestLogger(logger))
	r.Use(apimiddleware.Recoverer(logger))
	r.Use(httprate.LimitByIP(60, 1*time.Minute)) // Rate limit: 60 requests per minute per IP

	// Set up routes
	api.SetupRoutes(r, dbClient, authService, b2Service, wsHub, aiProcessor)

	// Start server
	addr := ":8080"
	logger.Info("Starting server", "address", addr)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", "error", err)
		}
	}()

	// Graceful shutdown
	waitForShutdown(server, logger)
}

func loadConfig() error {
	err := godotenv.Load(".env.local")
	if err != nil {
		return fmt.Errorf("error loading .env.local file: %v", err)
	}

	requiredEnvVars := []string{
		"CLERK_SECRET_KEY",
		"BACKBLAZE_ACCOUNT_ID",
		"BACKBLAZE_APPLICATION_KEY",
		"BACKBLAZE_BUCKET_ID",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%s environment variable is not set", envVar)
		}
	}
	return nil
}

func waitForShutdown(server *http.Server, logger *logger.Logger) {
	// ... (implementation of graceful shutdown)
}
