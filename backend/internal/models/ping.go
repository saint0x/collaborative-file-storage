package models

import (
	"time"

	"github.com/google/uuid"
)

type Ping struct {
	ID        uuid.UUID `json:"id"`
	ClientID  string    `json:"client_id"`
	Timestamp time.Time `json:"timestamp"`
}
