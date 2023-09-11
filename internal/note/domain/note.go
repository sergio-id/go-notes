package domain

import (
	"github.com/google/uuid"
	"time"
)

type Note struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	CategoryID *uuid.UUID `json:"category_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Pinned     bool       `json:"pinned"`
	Priority   int32      `json:"priority"`
	UpdatedAt  time.Time  `json:"update_at"`
	CreatedAt  time.Time  `json:"create_at"`
}
