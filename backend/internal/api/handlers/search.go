package handlers

import (
	"net/http"

	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func SearchFiles(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		// Implement file search logic
		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Search files", "query": query})
	}
}

func SearchFriends(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		// Implement friend search logic
		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Search friends", "query": query})
	}
}
