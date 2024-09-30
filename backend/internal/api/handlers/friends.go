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
)

func AddFriend(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var friendRequest struct {
			FriendID string `json:"friend_id"`
		}
		err = json.NewDecoder(r.Body).Decode(&friendRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		friend := models.Friend{
			ID:        uuid.New(),
			UserID:    userID,
			FriendID:  friendRequest.FriendID,
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = db.DB.Exec("INSERT INTO friends (id, user_id, friend_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			friend.ID, friend.UserID, friend.FriendID, friend.Status, friend.CreatedAt, friend.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to add friend", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(friend)
	}
}

func GetFriends(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.DB.Query("SELECT * FROM friends WHERE user_id = ? OR friend_id = ?", userID, userID)
		if err != nil {
			http.Error(w, "Failed to fetch friends", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var friends []models.Friend
		for rows.Next() {
			var f models.Friend
			err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
			if err != nil {
				http.Error(w, "Failed to scan friend", http.StatusInternalServerError)
				return
			}
			friends = append(friends, f)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(friends)
	}
}

func UpdateFriendStatus(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		friendshipID := chi.URLParam(r, "id")

		var updateData struct {
			Status string `json:"status"`
		}
		err = json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		_, err = db.DB.Exec("UPDATE friends SET status = ?, updated_at = ? WHERE id = ? AND (user_id = ? OR friend_id = ?)",
			updateData.Status, time.Now(), friendshipID, userID, userID)
		if err != nil {
			http.Error(w, "Failed to update friend status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Friend status updated successfully"})
	}
}

func RemoveFriend(db *db.SQLiteClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.GetUserIDFromContext(r.Context())
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		friendshipID := chi.URLParam(r, "id")

		_, err = db.DB.Exec("DELETE FROM friends WHERE id = ? AND (user_id = ? OR friend_id = ?)", friendshipID, userID, userID)
		if err != nil {
			http.Error(w, "Failed to remove friend", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Friend removed successfully"})
	}
}
