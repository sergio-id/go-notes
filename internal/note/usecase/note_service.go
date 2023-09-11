package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/sergio-id/go-notes/internal/note/domain"

	"github.com/google/wire"
)

type usecase struct {
	repo domain.NoteRepository
}

var _ domain.NoteUseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewNoteUseCase)

func NewNoteUseCase(
	repo domain.NoteRepository,
) domain.NoteUseCase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) ListNotes(ctx context.Context, arg *domain.ListNotesParams) ([]*domain.Note, error) {
	return u.repo.ListNotes(ctx, arg)
}

func (u *usecase) GetNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.Note, error) {
	return u.repo.GetNote(ctx, userID, id)
}

func (u *usecase) CreateNote(ctx context.Context, arg *domain.CreateNoteParams) (*domain.Note, error) {
	return u.repo.CreateNote(ctx, arg)
}

func (u *usecase) UpdateNote(ctx context.Context, arg *domain.UpdateNoteParams) (*domain.Note, error) {
	return u.repo.UpdateNote(ctx, arg)
}

func (u *usecase) DeleteNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	return u.repo.DeleteNote(ctx, userID, id)
}
