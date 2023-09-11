package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sergio-id/go-notes/internal/user/domain"

	"github.com/google/wire"
)

type usecase struct {
	repo domain.UserRepository
}

var _ domain.UserUseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUserUseCase)

func NewUserUseCase(
	repo domain.UserRepository,
) domain.UserUseCase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.repo.GetUser(ctx, id)
}

func (u *usecase) CreateUser(ctx context.Context, arg *domain.CreateUserParams) (*domain.User, error) {
	return u.repo.CreateUser(ctx, arg)
}

func (u *usecase) UpdateUser(ctx context.Context, arg *domain.UpdateUserParams) (*domain.User, error) {
	return u.repo.UpdateUser(ctx, arg)
}

func (u *usecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteUser(ctx, id)
}

func (u *usecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return u.repo.GetUserByEmail(ctx, email)
}
