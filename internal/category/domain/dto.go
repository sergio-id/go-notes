package domain

import "github.com/google/uuid"

type (
	ListCategoriesParams struct {
		UserID uuid.UUID
		Limit  *int32
		Offset *int32
	}

	CreateCategoryParams struct {
		UserID   uuid.UUID
		Title    string
		Pinned   bool
		Priority int32
	}

	UpdateCategoryParams struct {
		ID       uuid.UUID
		UserID   uuid.UUID
		Title    string
		Pinned   bool
		Priority int32
	}
)
