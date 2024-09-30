package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           uuid.UUID `json:"id"`
	UserID       string    `json:"user_id"`
	CollectionID uuid.UUID `json:"collection_id,omitempty"`
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
