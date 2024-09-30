package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/models"
	"github.com/saint0x/file-storage-app/backend/internal/services/ai"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/pkg/errors"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func OrganizeFiles(db *db.SQLiteClient, aiProcessor *ai.Processor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}

		var req struct {
			FileIDs []string `json:"file_ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondError(w, errors.BadRequest("Invalid request body"))
			return
		}

		files, err := db.GetFilesByIDs(req.FileIDs)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch files"))
			return
		}

		aiReq := ai.FileOrganizationRequest{Files: files}
		aiResp, err := aiProcessor.OrganizeFiles(r.Context(), aiReq)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to organize files"))
			return
		}

		// Parse userID string to UUID
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Invalid user ID"))
			return
		}

		// Create folders and update file associations
		for _, folder := range aiResp.Folders {
			newFolder := models.Folder{
				UserID: userUUID,
				Name:   folder.Name,
			}
			folderID, err := db.CreateFolder(newFolder)
			if err != nil {
				utils.RespondError(w, errors.InternalServerError("Failed to create folder"))
				return
			}

			for _, fileName := range folder.Files {
				if err := db.UpdateFileFolder(fileName, folderID); err != nil {
					utils.RespondError(w, errors.InternalServerError("Failed to update file folder"))
					return
				}
			}
		}

		utils.RespondJSON(w, http.StatusOK, aiResp)
	}
}
