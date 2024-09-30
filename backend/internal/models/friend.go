package models

import (
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	FriendID  string    `json:"friend_id"`
	Status    string    `json:"status"` // e.g., "pending", "accepted", "blocked"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
