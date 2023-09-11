package domain

import (
	"context"
)

type (
	SessionRepository interface {
		CreateSession(ctx context.Context, arg *CreateSessionParams) (*Session, error)
		GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
		DeleteByID(ctx context.Context, sessionID string) error
	}

	SessionUseCase interface {
		CreateSession(ctx context.Context, arg *CreateSessionParams) (*Session, error)
		GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
		DeleteByID(ctx context.Context, sessionID string) error
	}

	UserDomainService interface {
		GetUserByEmail(ctx context.Context, email string) (*UserDTO, error)
		CreateUser(ctx context.Context, arg *CreateUserParams) (*UserDTO, error)
	}
)
