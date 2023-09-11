package domain

import "github.com/google/uuid"

type (
	CreateUserParams struct {
		Email    string
		Password string
	}

	UpdateUserParams struct {
		ID        uuid.UUID
		FirstName string
		LastName  string
	}
)
