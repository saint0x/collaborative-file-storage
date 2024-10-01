package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/saint0x/file-storage-app/backend/internal/api"
	"github.com/saint0x/file-storage-app/backend/internal/api/handlers"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/models"
	"github.com/saint0x/file-storage-app/backend/internal/services/ai"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: Add user to database
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the uploads directory if it doesn't exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(filepath.Join("./uploads", header.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}

	// Initialize database
	dbClient, err := db.NewSQLiteClient(os.Getenv("SQLITE_DB_PATH"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbClient.Close()

	// Initialize B2 service
	b2Service, err := storage.NewB2Service(
		os.Getenv("BACKBLAZE_ACCOUNT_ID"),
		os.Getenv("BACKBLAZE_APPLICATION_KEY"),
		os.Getenv("BACKBLAZE_BUCKET_ID"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize B2 service: %v", err)
	}
	defer b2Service.Close()

	// Initialize Clerk service
	clerkService := auth.NewClerkService()
	clerkService.SetSecretKey(os.Getenv("CLERK_SECRET_KEY"))

	// Initialize WebSocket hub
	wsHub := websocket.NewHub(dbClient)
	go wsHub.Run()

	// Initialize AI processor
	aiProcessor := ai.NewProcessor(os.Getenv("OPENAI_API_KEY"))

	// Initialize router
	router := chi.NewRouter()

	// Initialize API routes
	api.SetupRoutes(router, dbClient, clerkService, b2Service, wsHub, aiProcessor)

	// Add routes for user creation and file upload
	router.Post("/users", createUser)
	router.Post("/upload", handlers.UploadFile(b2Service, dbClient))

	// Add health check route
	router.Get("/health", healthCheck)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
