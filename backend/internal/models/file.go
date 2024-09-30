package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	CollectionID *uuid.UUID `json:"collection_id,omitempty"`
	Key          string     `json:"key"`
	Name         string     `json:"name"`
	ContentType  string     `json:"content_type"`
	Size         int64      `json:"size"`
	UploadedAt   time.Time  `json:"uploaded_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
