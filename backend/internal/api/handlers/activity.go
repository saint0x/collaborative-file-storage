package handlers

import (
	"net/http"

	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
	"github.com/saint0x/file-storage-app/backend/pkg/errors"
	"github.com/saint0x/file-storage-app/backend/pkg/utils"
)

func GetRecentActivity(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			utils.RespondError(w, errors.Unauthorized("User not authenticated"))
			return
		}

		activity, err := db.GetRecentActivity(userID)
		if err != nil {
			utils.RespondError(w, errors.InternalServerError("Failed to fetch recent activity"))
			return
		}
		utils.RespondJSON(w, http.StatusOK, activity)
	}
}
