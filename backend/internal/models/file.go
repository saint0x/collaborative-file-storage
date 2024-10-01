package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Name         string
	ContentType  string
	Key          string
	Size         int64
	UploadedAt   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CollectionID uuid.UUID
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
