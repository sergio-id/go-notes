package domain

import "github.com/google/uuid"

type (
	ListNotesParams struct {
		UserID uuid.UUID
		Limit  *int32
		Offset *int32
	}

	CreateNoteParams struct {
		UserID     uuid.UUID
		CategoryID uuid.NullUUID
		Title      string
		Content    string
		Pinned     bool
		Priority   int32
	}

	UpdateNoteParams struct {
		ID       uuid.UUID
		UserID   uuid.UUID
		Title    string
		Content  string
		Pinned   bool
		Priority int32
	}
)
