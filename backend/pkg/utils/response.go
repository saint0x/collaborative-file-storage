package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// RespondError sends an error response
func RespondError(w http.ResponseWriter, err error) {
	// You can implement custom error handling here
	// For example, you can check for specific error types and set the status code accordingly
	statusCode := http.StatusInternalServerError

	// Example: Check for specific error types
	// if errors.Is(err, SomeCustomError) {
	//     statusCode = http.StatusBadRequest
	// }

	http.Error(w, err.Error(), statusCode)
}
