package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func GetSharedItems(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement get shared items logic
		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Get shared items"})
	}
}

func ShareItem(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement share item logic
		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Item shared successfully"})
	}
}

func UnshareItem(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemID := chi.URLParam(r, "id")
		// Implement unshare item logic
		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Item unshared successfully", "id": itemID})
	}
}
