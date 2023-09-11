package domain

import (
	"context"
	"github.com/google/uuid"
)

type (
	NoteRepository interface {
		ListNotes(ctx context.Context, arg *ListNotesParams) ([]*Note, error)
		GetNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*Note, error)
		CreateNote(ctx context.Context, arg *CreateNoteParams) (*Note, error)
		UpdateNote(ctx context.Context, arg *UpdateNoteParams) (*Note, error)
		DeleteNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	}

	NoteUseCase interface {
		ListNotes(ctx context.Context, arg *ListNotesParams) ([]*Note, error)
		GetNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*Note, error)
		CreateNote(ctx context.Context, arg *CreateNoteParams) (*Note, error)
		UpdateNote(ctx context.Context, arg *UpdateNoteParams) (*Note, error)
		DeleteNote(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	}

	AuthDomainService interface {
		GetUserIDByToken(ctx context.Context, token string) (uuid.UUID, error)
	}
)
