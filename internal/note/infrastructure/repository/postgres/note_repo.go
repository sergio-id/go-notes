package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sergio-id/go-notes/internal/note/domain"
	postgresql "github.com/sergio-id/go-notes/internal/note/infrastructure/repository/postgres/sqlc"
	"github.com/sergio-id/go-notes/pkg/postgres"
	"time"
)

const (
	_defaultLimit  = 50
	_defaultOffset = 0
)

type noteRepo struct {
	pg postgres.DBEngine
}

var _ domain.NoteRepository = (*noteRepo)(nil)

var RepositorySet = wire.NewSet(NewNoteRepo)

func NewNoteRepo(pg postgres.DBEngine) domain.NoteRepository {
	return &noteRepo{pg: pg}
}

func (r *noteRepo) ListNotes(ctx context.Context, arg *domain.ListNotesParams) ([]*domain.Note, error) {
	q := postgresql.New(r.pg.GetDB())

	limit := int32(_defaultLimit)
	if arg.Limit == nil {
		limit = 10
	}
	offset := int32(_defaultOffset)
	if arg.Offset == nil {
		offset = 0
	}

	fetchedNotes, err := q.ListNotes(ctx, postgresql.ListNotesParams{
		UserID: arg.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.ListNotes")
	}

	notes := make([]*domain.Note, 0, len(fetchedNotes))
	for _, fetchedNote := range fetchedNotes {
		notes = append(notes, &domain.Note{
			ID:         fetchedNote.ID,
			UserID:     fetchedNote.UserID,
			CategoryID: &fetchedNote.CategoryID.UUID,
			Title:      fetchedNote.Title,
			Content:    fetchedNote.Content,
			Pinned:     fetchedNote.Pinned,
			Priority:   fetchedNote.Priority,
			UpdatedAt:  fetchedNote.UpdatedAt,
			CreatedAt:  fetchedNote.CreatedAt,
		})
	}

	return notes, nil
}

func (r *noteRepo) GetNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.Note, error) {
	q := postgresql.New(r.pg.GetDB())

	fetchedNote, err := q.GetNote(ctx, postgresql.GetNoteParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetNote")
	}

	return &domain.Note{
		ID:         fetchedNote.ID,
		UserID:     fetchedNote.UserID,
		CategoryID: &fetchedNote.CategoryID.UUID,
		Title:      fetchedNote.Title,
		Content:    fetchedNote.Content,
		Pinned:     fetchedNote.Pinned,
		Priority:   fetchedNote.Priority,
		UpdatedAt:  fetchedNote.UpdatedAt,
		CreatedAt:  fetchedNote.CreatedAt,
	}, nil
}

func (r *noteRepo) CreateNote(ctx context.Context, arg *domain.CreateNoteParams) (*domain.Note, error) {
	q := postgresql.New(r.pg.GetDB())

	createdNote, err := q.CreateNote(ctx, postgresql.CreateNoteParams{
		ID:         uuid.New(),
		UserID:     arg.UserID,
		CategoryID: arg.CategoryID,
		Title:      arg.Title,
		Content:    arg.Content,
		Pinned:     arg.Pinned,
		Priority:   arg.Priority,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.CreateNote")
	}

	note := domain.Note{
		ID:         createdNote.ID,
		UserID:     createdNote.UserID,
		CategoryID: nil,
		Title:      createdNote.Title,
		Content:    createdNote.Content,
		Pinned:     createdNote.Pinned,
		Priority:   createdNote.Priority,
		UpdatedAt:  createdNote.UpdatedAt,
		CreatedAt:  createdNote.CreatedAt,
	}

	if createdNote.CategoryID.Valid {
		note.CategoryID = &createdNote.CategoryID.UUID
	}

	return &note, nil
}

func (r *noteRepo) UpdateNote(ctx context.Context, arg *domain.UpdateNoteParams) (*domain.Note, error) {
	q := postgresql.New(r.pg.GetDB())

	updatedNote, err := q.UpdateNote(ctx, postgresql.UpdateNoteParams{
		ID:        arg.ID,
		Title:     arg.Title,
		Content:   arg.Content,
		Pinned:    arg.Pinned,
		Priority:  arg.Priority,
		UpdatedAt: time.Now(),
		UserID:    arg.UserID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "querier.UpdateNote")
	}

	return &domain.Note{
		ID:         updatedNote.ID,
		UserID:     updatedNote.UserID,
		CategoryID: &updatedNote.CategoryID.UUID,
		Title:      updatedNote.Title,
		Content:    updatedNote.Content,
		Pinned:     updatedNote.Pinned,
		Priority:   updatedNote.Priority,
		UpdatedAt:  updatedNote.UpdatedAt,
		CreatedAt:  updatedNote.CreatedAt,
	}, nil
}

func (r *noteRepo) DeleteNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	q := postgresql.New(r.pg.GetDB())

	err := q.DeleteNote(ctx, postgresql.DeleteNoteParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return errors.Wrap(err, "querier.DeleteNote")
	}

	return nil
}
