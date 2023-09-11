package domain

import (
	"github.com/google/uuid"
	"time"
)

type (
	UserDTO struct {
		ID       uuid.UUID `json:"id"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
	}

	CreateUserParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	CreateSessionParams struct {
		UserID   uuid.UUID     `json:"user_id"`
		Duration time.Duration `json:"duration"`
	}
)
