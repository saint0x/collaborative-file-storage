package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/saint0x/file-storage-app/backend/internal/api/handlers"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/health", handlers.HealthCheck)

	// Add more routes here

	return r
}
