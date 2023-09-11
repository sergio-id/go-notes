package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/internal/user/domain"
	postgresql "github.com/sergio-id/go-notes/internal/user/infrastructure/repository/postgres/sqlc"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"time"
)

type userRepo struct {
	pg postgres.DBEngine
}

var _ domain.UserRepository = (*userRepo)(nil)

var RepositorySet = wire.NewSet(NewUserRepo)

func NewUserRepo(pg postgres.DBEngine) domain.UserRepository {
	return &userRepo{pg: pg}
}

func (r *userRepo) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	q := postgresql.New(r.pg.GetDB())

	fetchedUser, err := q.GetUser(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetUser")
	}

	return &domain.User{
		ID:        fetchedUser.ID,
		Email:     fetchedUser.Email,
		Password:  fetchedUser.Password,
		FirstName: fetchedUser.FirstName,
		LastName:  fetchedUser.LastName,
		UpdatedAt: fetchedUser.UpdatedAt,
		CreatedAt: fetchedUser.CreatedAt,
	}, nil
}

func (r *userRepo) CreateUser(ctx context.Context, arg *domain.CreateUserParams) (*domain.User, error) {
	q := postgresql.New(r.pg.GetDB())

	createdUser, err := q.CreateUser(ctx, postgresql.CreateUserParams{
		ID:        uuid.New(),
		Email:     arg.Email,
		Password:  arg.Password,
		FirstName: "",
		LastName:  "",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.CreateUser")
	}

	return &domain.User{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		Password:  createdUser.Password,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		UpdatedAt: createdUser.UpdatedAt,
		CreatedAt: createdUser.CreatedAt,
	}, nil
}

func (r *userRepo) UpdateUser(ctx context.Context, arg *domain.UpdateUserParams) (*domain.User, error) {
	q := postgresql.New(r.pg.GetDB())

	updatedUser, err := q.UpdateUser(ctx, postgresql.UpdateUserParams{
		ID:        arg.ID,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.UpdateUser")
	}

	return &domain.User{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Password:  updatedUser.Password,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		UpdatedAt: updatedUser.UpdatedAt,
		CreatedAt: updatedUser.CreatedAt,
	}, nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	q := postgresql.New(r.pg.GetDB())

	err := q.DeleteUser(ctx, id)
	if err != nil {
		return errors.Wrap(err, "querier.DeleteUser")
	}

	return nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	q := postgresql.New(r.pg.GetDB())

	fetchedUser, err := q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetUserByEmail")
	}

	return &domain.User{
		ID:        fetchedUser.ID,
		Email:     fetchedUser.Email,
		Password:  fetchedUser.Password,
		FirstName: fetchedUser.FirstName,
		LastName:  fetchedUser.LastName,
		UpdatedAt: fetchedUser.UpdatedAt,
		CreatedAt: fetchedUser.CreatedAt,
	}, nil
}
