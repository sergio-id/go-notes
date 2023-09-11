package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/internal/category/domain"
	postgresql "github.com/sergio-id/go-notes/internal/category/infrastructure/repository/postgres/sqlc"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"time"
)

const (
	_defaultLimit  = 50
	_defaultOffset = 0
)

type categoryRepo struct {
	pg postgres.DBEngine
}

var _ domain.CategoryRepository = (*categoryRepo)(nil)

var RepositorySet = wire.NewSet(NewCategoryRepo)

func NewCategoryRepo(pg postgres.DBEngine) domain.CategoryRepository {
	return &categoryRepo{pg: pg}
}

func (r *categoryRepo) ListCategories(ctx context.Context, arg *domain.ListCategoriesParams) ([]*domain.Category, error) {
	q := postgresql.New(r.pg.GetDB())

	limit := int32(_defaultLimit)
	if arg.Limit != nil {
		limit = *arg.Limit
	}

	offset := int32(_defaultOffset)
	if arg.Offset != nil {
		offset = *arg.Offset
	}

	fetchedCategories, err := q.ListCategories(ctx, postgresql.ListCategoriesParams{
		UserID: arg.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.ListCategories")
	}

	categories := make([]*domain.Category, len(fetchedCategories))
	for i, fetchedCategory := range fetchedCategories {
		categories[i] = &domain.Category{
			ID:        fetchedCategory.ID,
			UserID:    fetchedCategory.UserID,
			Title:     fetchedCategory.Title,
			Pinned:    fetchedCategory.Pinned,
			Priority:  fetchedCategory.Priority,
			UpdatedAt: fetchedCategory.UpdatedAt,
			CreatedAt: fetchedCategory.CreatedAt,
		}
	}

	return categories, nil
}

func (r *categoryRepo) GetCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.Category, error) {
	q := postgresql.New(r.pg.GetDB())

	fetchedCategory, err := q.GetCategory(ctx, postgresql.GetCategoryParams{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetCategory")
	}

	return &domain.Category{
		ID:        fetchedCategory.ID,
		UserID:    fetchedCategory.UserID,
		Title:     fetchedCategory.Title,
		Pinned:    fetchedCategory.Pinned,
		Priority:  fetchedCategory.Priority,
		UpdatedAt: fetchedCategory.UpdatedAt,
		CreatedAt: fetchedCategory.CreatedAt,
	}, nil
}

func (r *categoryRepo) CreateCategory(ctx context.Context, arg *domain.CreateCategoryParams) (*domain.Category, error) {
	q := postgresql.New(r.pg.GetDB())

	createdCategory, err := q.CreateCategory(ctx, postgresql.CreateCategoryParams{
		ID:        uuid.New(),
		UserID:    arg.UserID,
		Title:     arg.Title,
		Pinned:    arg.Pinned,
		Priority:  arg.Priority,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.CreateCategory")
	}

	return &domain.Category{
		ID:        createdCategory.ID,
		UserID:    createdCategory.UserID,
		Title:     createdCategory.Title,
		Pinned:    createdCategory.Pinned,
		Priority:  createdCategory.Priority,
		UpdatedAt: createdCategory.UpdatedAt,
		CreatedAt: createdCategory.CreatedAt,
	}, nil
}

func (r *categoryRepo) UpdateCategory(ctx context.Context, arg *domain.UpdateCategoryParams) (*domain.Category, error) {
	q := postgresql.New(r.pg.GetDB())

	updatedCategory, err := q.UpdateCategory(ctx, postgresql.UpdateCategoryParams{
		ID:        arg.ID,
		UserID:    arg.UserID,
		Title:     arg.Title,
		Pinned:    arg.Pinned,
		Priority:  arg.Priority,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.UpdateCategory")
	}

	return &domain.Category{
		ID:        updatedCategory.ID,
		UserID:    updatedCategory.UserID,
		Title:     updatedCategory.Title,
		Pinned:    updatedCategory.Pinned,
		Priority:  updatedCategory.Priority,
		UpdatedAt: updatedCategory.UpdatedAt,
		CreatedAt: updatedCategory.CreatedAt,
	}, nil
}

func (r *categoryRepo) DeleteCategory(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	q := postgresql.New(r.pg.GetDB())

	err := q.DeleteCategory(ctx, postgresql.DeleteCategoryParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return errors.Wrap(err, "querier.DeleteCategory")
	}

	return nil
}
