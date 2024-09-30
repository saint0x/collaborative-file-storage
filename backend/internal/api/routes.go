package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/saint0x/file-storage-app/backend/internal/api/handlers"
	apimiddleware "github.com/saint0x/file-storage-app/backend/internal/api/middleware"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/services/ai"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
)

func SetupRoutes(
	db *db.SQLiteClient,
	authService *auth.ClerkService,
	storageService *storage.R2Service,
	wsHub *websocket.Hub,
	aiProcessor *ai.Processor,
) http.Handler {
	r := chi.NewRouter()

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Adjust this to your needs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(corsMiddleware.Handler)
	r.Use(apimiddleware.Logging(wsHub)) // Add the new logging middleware

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth/webhook", handlers.ClerkWebhook(authService, db))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(apimiddleware.Auth(authService))

		// User routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", handlers.GetUsers(db))
			r.Get("/{id}", handlers.GetUser(db))
			r.Put("/{id}", handlers.UpdateUser(db))
		})

		// Collection routes
		r.Route("/collections", func(r chi.Router) {
			r.Get("/", handlers.GetCollections(db))
			r.Post("/", handlers.CreateCollection(db))
			r.Put("/{id}", handlers.UpdateCollection(db))
			r.Delete("/{id}", handlers.DeleteCollection(db))
		})

		// File routes
		r.Route("/files", func(r chi.Router) {
			r.Get("/", handlers.GetFiles(db))
			r.Post("/", handlers.UploadFile(db, storageService))
			r.Delete("/{id}", handlers.DeleteFile(db, storageService))
		})

		// Friend routes
		r.Route("/friends", func(r chi.Router) {
			r.Get("/", handlers.GetFriends(db))
			r.Post("/", handlers.AddFriend(db))
			r.Put("/{id}", handlers.UpdateFriendStatus(db))
			r.Delete("/{id}", handlers.RemoveFriend(db))
		})

		// WebSocket route
		r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			websocket.ServeWs(wsHub, w, r)
		})

		// Organize files route
		r.Post("/organize-files", handlers.OrganizeFiles(db, aiProcessor))
	})

	return r
}
