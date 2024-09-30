package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/models"
	"github.com/saint0x/file-storage-app/backend/internal/services/auth"
)

func ClerkWebhook(authService *auth.ClerkService, db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var webhookData struct {
			Data struct {
				ID        string `json:"id"`
				Email     string `json:"email_address"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
			} `json:"data"`
			Type string `json:"type"`
		}

		err := json.NewDecoder(r.Body).Decode(&webhookData)
		if err != nil {
			http.Error(w, "Invalid webhook data", http.StatusBadRequest)
			return
		}

		switch webhookData.Type {
		case "user.created":
			newUser := models.User{
				ID:        uuid.New(),
				ClerkID:   webhookData.Data.ID,
				Email:     webhookData.Data.Email,
				Username:  webhookData.Data.Username,
				FirstName: webhookData.Data.FirstName,
				LastName:  webhookData.Data.LastName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			_, err = db.DB.Exec("INSERT INTO users (id, clerk_id, email, username, first_name, last_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
				newUser.ID, newUser.ClerkID, newUser.Email, newUser.Username, newUser.FirstName, newUser.LastName, newUser.CreatedAt, newUser.UpdatedAt)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}

		case "user.updated":
			_, err = db.DB.Exec("UPDATE users SET email = ?, username = ?, first_name = ?, last_name = ?, updated_at = ? WHERE clerk_id = ?",
				webhookData.Data.Email, webhookData.Data.Username, webhookData.Data.FirstName, webhookData.Data.LastName, time.Now(), webhookData.Data.ID)
			if err != nil {
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return
			}

		case "user.deleted":
			_, err = db.DB.Exec("DELETE FROM users WHERE clerk_id = ?", webhookData.Data.ID)
			if err != nil {
				http.Error(w, "Failed to delete user", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Webhook processed successfully"})
	}
}
