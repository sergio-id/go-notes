package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sergio-id/go-notes/internal/category/domain"

	"github.com/google/wire"
)

type usecase struct {
	repo domain.CategoryRepository
}

var _ domain.CategoryUseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewCategoryUseCase)

func NewCategoryUseCase(
	repo domain.CategoryRepository,
) domain.CategoryUseCase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) ListCategories(ctx context.Context, arg *domain.ListCategoriesParams) ([]*domain.Category, error) {
	return u.repo.ListCategories(ctx, arg)
}

func (u *usecase) GetCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.Category, error) {
	return u.repo.GetCategory(ctx, userID, id)
}

func (u *usecase) CreateCategory(ctx context.Context, arg *domain.CreateCategoryParams) (*domain.Category, error) {
	return u.repo.CreateCategory(ctx, arg)
}

func (u *usecase) UpdateCategory(ctx context.Context, arg *domain.UpdateCategoryParams) (*domain.Category, error) {
	return u.repo.UpdateCategory(ctx, arg)
}

func (u *usecase) DeleteCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return u.repo.DeleteCategory(ctx, userID, id)
}
