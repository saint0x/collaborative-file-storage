package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/models"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/pkg/errors"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func UploadFile(b2Service *storage.B2Service, dbClient *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, fmt.Errorf("unauthorized: %w", err))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			utils.RespondError(w, fmt.Errorf("failed to get file from form: %w", err))
			return
		}
		defer file.Close()

		// Generate a unique key for the file
		fileID := uuid.New().String()
		key := fmt.Sprintf("%s/%s_%s", userID, fileID, header.Filename)

		// Upload the file to B2
		err = b2Service.UploadFile(r.Context(), key, file)
		if err != nil {
			utils.RespondError(w, fmt.Errorf("failed to upload file to B2: %w", err))
			return
		}

		// Create a new File record
		newFile := models.File{
			ID:          uuid.MustParse(fileID),
			UserID:      uuid.MustParse(userID),
			Key:         key,
			Name:        header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			Size:        header.Size,
			UploadedAt:  time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Save the file metadata to the database
		err = dbClient.CreateFile(newFile)
		if err != nil {
			utils.RespondError(w, fmt.Errorf("failed to save file metadata: %w", err))
			return
		}

		utils.RespondJSON(w, http.StatusCreated, newFile)
	}
}

func GetFiles(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}

		pagination, err := utils.NewPaginationFromRequest(r.URL.Query().Get("page"), r.URL.Query().Get("page_size"))
		if err != nil {
			utils.RespondError(w, err)
			return
		}

		rows, err := db.DB.Query("SELECT * FROM files WHERE user_id = ? LIMIT ? OFFSET ?", userID, pagination.PageSize, pagination.CalculateOffset())
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch files"))
			return
		}
		defer rows.Close()

		var files []models.File
		for rows.Next() {
			var f models.File
			var collectionID sql.NullString
			err := rows.Scan(&f.ID, &f.UserID, &collectionID, &f.Key, &f.Name, &f.Size, &f.ContentType, &f.UploadedAt, &f.CreatedAt, &f.UpdatedAt)
			if err != nil {
				utils.RespondError(w, errors.InternalServerError("Failed to scan file"))
				return
			}
			if collectionID.Valid {
				collUUID, _ := uuid.Parse(collectionID.String)
				f.CollectionID = uuid.NullUUID{UUID: collUUID, Valid: true}
			}
			files = append(files, f)
		}

		// Get total count for pagination
		var totalCount int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM files WHERE user_id = ?", userID).Scan(&totalCount)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to get total file count"))
			return
		}

		paginationInfo := utils.CalculatePagination(totalCount, pagination.Page, pagination.PageSize)

		response := map[string]interface{}{
			"files":      files,
			"pagination": paginationInfo,
		}

		utils.RespondJSON(w, http.StatusOK, response)
	}
}

func DeleteFile(db *db.SQLiteClient, storageService *storage.B2Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}
		fileID := chi.URLParam(r, "id")

		var key string
		err = db.DB.QueryRow("SELECT key FROM files WHERE id = ? AND user_id = ?", fileID, userID).Scan(&key)
		if err != nil {
			utils.RespondError(w, errors.NotFound("File not found or not owned by user"))
			return
		}

		err = storageService.DeleteFile(r.Context(), key)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to delete file from storage"))
			return
		}

		_, err = db.DB.Exec("DELETE FROM files WHERE id = ?", fileID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to delete file metadata"))
			return
		}

		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "File deleted successfully"})
	}
}

func GetFileCategories(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := db.GetFileCategories()
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch file categories"))
			return
		}
		utils.RespondJSON(w, http.StatusOK, categories)
	}
}

func GetFilesByCategory(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categoryName := chi.URLParam(r, "categoryName")
		files, err := db.GetFilesByCategory(categoryName)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch files by category"))
			return
		}
		utils.RespondJSON(w, http.StatusOK, files)
	}
}

func GetFileDetails(db *db.SQLiteClient, storageService *storage.B2Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileID := chi.URLParam(r, "id")
		details, err := db.GetFileDetails(fileID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch file details"))
			return
		}
		// Add additional details from storage service if needed
		utils.RespondJSON(w, http.StatusOK, details)
	}
}

func ShareFileWithFriends(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileID := chi.URLParam(r, "id")
		var req struct {
			FriendIDs []string `json:"friend_ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, errors.BadRequest("Invalid request body"))
			return
		}
		if err := db.ShareFileWithFriends(fileID, req.FriendIDs); err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to share file with friends"))
			return
		}
		utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "File shared successfully"})
	}
}

func GetSharedWithMeFiles(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := auth.GetUserIDFromContext(r.Context())
		files, err := db.GetSharedWithMeFiles(userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch shared files"))
			return
		}
		utils.RespondJSON(w, http.StatusOK, files)
	}
}

func GetOrganizedFileStructure(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := auth.GetUserIDFromContext(r.Context())
		structure, err := db.GetOrganizedFileStructure(userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch organized file structure"))
			return
		}
		utils.RespondJSON(w, http.StatusOK, structure)
	}
}

func UpdateFile(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation for updating file metadata
	}
}
