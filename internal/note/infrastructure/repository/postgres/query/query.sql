-- name: ListNotes :many
SELECT *
FROM note.note
WHERE user_id = $1
ORDER BY pinned DESC, priority DESC, title ASC LIMIT $2
OFFSET $3;

-- name: GetNote :one
SELECT *
FROM note.note
WHERE id = $1
  AND user_id = $2;

-- name: CreateNote :one
INSERT INTO note.note (id, user_id, category_id, title, content, pinned, priority, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateNote :one
UPDATE note.note
SET title      = $2,
    content    = $3,
    pinned     = $4,
    priority   = $5,
    updated_at = $6
WHERE id = $1 AND user_id = $7 RETURNING *;

-- name: DeleteNote :exec
DELETE
FROM note.note
WHERE id = $1
  AND user_id = $2;
