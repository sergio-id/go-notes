package domain

import (
	"github.com/google/uuid"
	"time"
)

type Category struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Pinned    bool      `json:"pinned"`
	Priority  int32     `json:"priority"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
