package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/models"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/pkg/errors"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func CreateCollection(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}

		var collection models.Collection
		err = json.NewDecoder(r.Body).Decode(&collection)
		if err != nil {
			utils.RespondError(w, errors.BadRequest("Invalid request body"))
			return
		}

		collection.ID = uuid.New()
		collection.UserID = userID
		collection.CreatedAt = time.Now()
		collection.UpdatedAt = time.Now()

		_, err = db.DB.Exec("INSERT INTO collections (id, user_id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			collection.ID, collection.UserID, collection.Name, collection.Description, collection.CreatedAt, collection.UpdatedAt)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to create collection"))
			return
		}

		utils.RespondJSON(w, http.StatusCreated, collection)
	}
}

func GetCollections(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}

		rows, err := db.DB.Query("SELECT * FROM collections WHERE user_id = ?", userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch collections"))
			return
		}
		defer rows.Close()

		var collections []models.Collection
		for rows.Next() {
			var c models.Collection
			err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
			if err != nil {
				utils.RespondError(w, errors.InternalServerError("Failed to scan collection"))
				return
			}
			collections = append(collections, c)
		}

		utils.RespondJSON(w, http.StatusOK, collections)
	}
}

func UpdateCollection(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}
		collectionID := chi.URLParam(r, "id")

		var updateData models.Collection
		err = json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			utils.RespondError(w, errors.BadRequest("Invalid request body"))
			return
		}

		_, err = db.DB.Exec("UPDATE collections SET name = ?, description = ?, updated_at = ? WHERE id = ? AND user_id = ?",
			updateData.Name, updateData.Description, time.Now(), collectionID, userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to update collection"))
			return
		}

		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Collection updated successfully"})
	}
}

func DeleteCollection(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}
		collectionID := chi.URLParam(r, "id")

		_, err = db.DB.Exec("DELETE FROM collections WHERE id = ? AND user_id = ?", collectionID, userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to delete collection"))
			return
		}

		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Collection deleted successfully"})
	}
}
