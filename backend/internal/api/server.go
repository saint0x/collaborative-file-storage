package api

import (
	"fmt"
	"net/http"

	"github.com/saint0x/file-storage-app/backend/internal/api/routes"
)

func StartServer() error {
	router := routes.SetupRoutes()

	port := ":8080" // You can make this configurable if needed
	fmt.Printf("Server starting on port %s\n", port)

	return http.ListenAndServe(port, router)
}
