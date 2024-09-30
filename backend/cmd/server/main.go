package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/saint0x/file-storage-app/backend/internal/api"
	"github.com/saint0x/file-storage-app/backend/internal/config"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/services/ai"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
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

	// Initialize router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize services
	authService, err := auth.NewClerkService()
	if err != nil {
		log.Fatalf("Failed to initialize Clerk service: %v", err)
	}
	storageService, err := storage.NewR2Service(cfg.R2AccountID, cfg.R2AccessKeyID, cfg.R2SecretAccessKey, cfg.R2BucketName)
	if err != nil {
		log.Fatalf("Failed to initialize R2 service: %v", err)
	}
	wsHub := websocket.NewHub(database)
	aiProcessor := ai.NewProcessor(cfg.OpenAIAPIKey)

	// Set up routes
	handler := api.SetupRoutes(database, authService, storageService, wsHub, aiProcessor)
	r.Mount("/", handler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Initialize WebSocket hub
	go wsHub.Run()
}
