package domain

import (
	"context"
	"github.com/google/uuid"
)

type (
	UserRepository interface {
		GetUser(ctx context.Context, id uuid.UUID) (*User, error)
		CreateUser(ctx context.Context, arg *CreateUserParams) (*User, error)
		UpdateUser(ctx context.Context, arg *UpdateUserParams) (*User, error)
		DeleteUser(ctx context.Context, id uuid.UUID) error
		GetUserByEmail(ctx context.Context, email string) (*User, error)
	}

	UserUseCase interface {
		GetUser(ctx context.Context, id uuid.UUID) (*User, error)
		CreateUser(ctx context.Context, arg *CreateUserParams) (*User, error)
		UpdateUser(ctx context.Context, arg *UpdateUserParams) (*User, error)
		DeleteUser(ctx context.Context, id uuid.UUID) error
		GetUserByEmail(ctx context.Context, email string) (*User, error)
	}

	AuthDomainService interface {
		GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
	}
)
