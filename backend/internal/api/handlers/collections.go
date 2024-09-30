package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/yourproject/internal/db"
	"github.com/yourusername/yourproject/internal/models"
	"github.com/yourusername/yourproject/internal/services/auth"
)

func GetCollections(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := auth.GetUserIDFromContext(r.Context())

		rows, err := db.DB.Query("SELECT * FROM collections WHERE user_id = ?", userID)
		if err != nil {
			http.Error(w, "Failed to fetch collections", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var collections []models.Collection
		for rows.Next() {
			var c models.Collection
			err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				http.Error(w, "Failed to scan collection", http.StatusInternalServerError)
				return
			}
			collections = append(collections, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(collections)
	}
}
