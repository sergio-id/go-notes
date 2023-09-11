// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: query.sql

package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createNote = `-- name: CreateNote :one
INSERT INTO note.note (id, user_id, category_id, title, content, pinned, priority, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, user_id, category_id, title, content, pinned, priority, updated_at, created_at
`

type CreateNoteParams struct {
	ID         uuid.UUID     `json:"id"`
	UserID     uuid.UUID     `json:"user_id"`
	CategoryID uuid.NullUUID `json:"category_id"`
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	Pinned     bool          `json:"pinned"`
	Priority   int32         `json:"priority"`
	UpdatedAt  time.Time     `json:"updated_at"`
	CreatedAt  time.Time     `json:"created_at"`
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (NoteNote, error) {
	row := q.db.QueryRowContext(ctx, createNote,
		arg.ID,
		arg.UserID,
		arg.CategoryID,
		arg.Title,
		arg.Content,
		arg.Pinned,
		arg.Priority,
		arg.UpdatedAt,
		arg.CreatedAt,
	)
	var i NoteNote
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CategoryID,
		&i.Title,
		&i.Content,
		&i.Pinned,
		&i.Priority,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE
FROM note.note
WHERE id = $1
  AND user_id = $2
`

type DeleteNoteParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteNote(ctx context.Context, arg DeleteNoteParams) error {
	_, err := q.db.ExecContext(ctx, deleteNote, arg.ID, arg.UserID)
	return err
}

const getNote = `-- name: GetNote :one
SELECT id, user_id, category_id, title, content, pinned, priority, updated_at, created_at
FROM note.note
WHERE id = $1
  AND user_id = $2
`

type GetNoteParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetNote(ctx context.Context, arg GetNoteParams) (NoteNote, error) {
	row := q.db.QueryRowContext(ctx, getNote, arg.ID, arg.UserID)
	var i NoteNote
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CategoryID,
		&i.Title,
		&i.Content,
		&i.Pinned,
		&i.Priority,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listNotes = `-- name: ListNotes :many
SELECT id, user_id, category_id, title, content, pinned, priority, updated_at, created_at
FROM note.note
WHERE user_id = $1
ORDER BY pinned DESC, priority DESC, title ASC LIMIT $2
OFFSET $3
`

type ListNotesParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListNotes(ctx context.Context, arg ListNotesParams) ([]NoteNote, error) {
	rows, err := q.db.QueryContext(ctx, listNotes, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NoteNote
	for rows.Next() {
		var i NoteNote
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CategoryID,
			&i.Title,
			&i.Content,
			&i.Pinned,
			&i.Priority,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateNote = `-- name: UpdateNote :one
UPDATE note.note
SET title      = $2,
    content    = $3,
    pinned     = $4,
    priority   = $5,
    updated_at = $6
WHERE id = $1 AND user_id = $7 RETURNING id, user_id, category_id, title, content, pinned, priority, updated_at, created_at
`

type UpdateNoteParams struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Pinned    bool      `json:"pinned"`
	Priority  int32     `json:"priority"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

func (q *Queries) UpdateNote(ctx context.Context, arg UpdateNoteParams) (NoteNote, error) {
	row := q.db.QueryRowContext(ctx, updateNote,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.Pinned,
		arg.Priority,
		arg.UpdatedAt,
		arg.UserID,
	)
	var i NoteNote
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CategoryID,
		&i.Title,
		&i.Content,
		&i.Pinned,
		&i.Priority,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}