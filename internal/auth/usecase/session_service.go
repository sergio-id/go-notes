package usecase

import (
	"context"
	"github.com/google/wire"
	"github.com/sergio-id/go-notes/internal/auth/domain"
)

type usecase struct {
	repo domain.SessionRepository
}

var _ domain.SessionUseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewSessionUseCase)

func NewSessionUseCase(
	repo domain.SessionRepository,
) domain.SessionUseCase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateSession(ctx context.Context, arg *domain.CreateSessionParams) (*domain.Session, error) {
	return u.repo.CreateSession(ctx, arg)
}

func (u *usecase) GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	return u.repo.GetSessionByID(ctx, sessionID)
}

func (u *usecase) DeleteByID(ctx context.Context, sessionID string) error {
	return u.repo.DeleteByID(ctx, sessionID)
}
