package domain

import (
	"context"
	"github.com/google/uuid"
)

type (
	CategoryRepository interface {
		ListCategories(ctx context.Context, arg *ListCategoriesParams) ([]*Category, error)
		GetCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*Category, error)
		CreateCategory(ctx context.Context, arg *CreateCategoryParams) (*Category, error)
		UpdateCategory(ctx context.Context, arg *UpdateCategoryParams) (*Category, error)
		DeleteCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	}

	CategoryUseCase interface {
		ListCategories(ctx context.Context, arg *ListCategoriesParams) ([]*Category, error)
		GetCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*Category, error)
		CreateCategory(ctx context.Context, arg *CreateCategoryParams) (*Category, error)
		UpdateCategory(ctx context.Context, arg *UpdateCategoryParams) (*Category, error)
		DeleteCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	}

	AuthDomainService interface {
		GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
	}
)
