package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/saint0x/file-storage-app/backend/internal/api/handlers"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/services/ai"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
)

func SetupRoutes(
	r *chi.Mux,
	db *db.SQLiteClient,
	authService *auth.ClerkService,
	storageService *storage.B2Service,
	wsHub *websocket.Hub,
	aiProcessor *ai.Processor,
) http.Handler {
	// ... (existing routes)

	// Sharing routes
	r.Route("/sharing", func(r chi.Router) {
		r.Get("/", handlers.GetSharedItems(db))
		r.Post("/", handlers.ShareItem(db))
		r.Delete("/{id}", handlers.UnshareItem(db))
	})

	// ... (existing routes)

	// Search routes
	r.Route("/search", func(r chi.Router) {
		r.Get("/files", handlers.SearchFiles(db))
		r.Get("/friends", handlers.SearchFriends(db))
	})

	// ... (rest of the function)

	return r
}
