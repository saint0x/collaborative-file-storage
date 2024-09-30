package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yourusername/yourproject/internal/api/handlers"
	"github.com/yourusername/yourproject/internal/api/middleware"
	"github.com/yourusername/yourproject/internal/db"
	"github.com/yourusername/yourproject/internal/services/ai"
	"github.com/yourusername/yourproject/internal/services/auth"
	"github.com/yourusername/yourproject/internal/services/storage"
	"github.com/yourusername/yourproject/internal/services/websocket"
)

func SetupRoutes(
	db *db.SQLiteClient,
	authService *auth.ClerkService,
	storageService *storage.R2Service,
	aiService *ai.Processor,
	wsHub *websocket.Hub,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(apimiddleware.Cors)

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
			r.Get("/{id}", handlers.GetCollection(db))
			r.Put("/{id}", handlers.UpdateCollection(db))
			r.Delete("/{id}", handlers.DeleteCollection(db))
		})

		// File routes
		r.Route("/files", func(r chi.Router) {
			r.Get("/", handlers.GetFiles(db))
			r.Post("/", handlers.UploadFile(db, storageService))
			r.Get("/{id}", handlers.GetFile(db))
			r.Delete("/{id}", handlers.DeleteFile(db, storageService))
		})

		// AI routes
		r.Route("/ai", func(r chi.Router) {
			r.Post("/process", handlers.ProcessAI(aiService))
		})

		// WebSocket route
		r.Get("/ws", handlers.ServeWs(wsHub))
	})

	return r
}
