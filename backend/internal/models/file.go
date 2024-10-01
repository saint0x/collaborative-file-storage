package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	FolderID     uuid.NullUUID // Add this field
	CollectionID uuid.NullUUID
	Key          string
	Name         string
	ContentType  string
	Size         int64
	UploadedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	B2FileID     string // Add this field
}

type FileDetails struct {
	File
	Size           int64
	UploadedAt     time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Collection     string
	Folder         string
	CollectionName string // Add this field
	FolderName     string // Add this field
}

type FileStructure struct {
	Folders map[string]Folder
	Files   map[string][]File
}
