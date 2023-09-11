-- name: ListCategories :many
SELECT *
FROM category.category
WHERE user_id = $1
ORDER BY pinned DESC, priority DESC, title ASC
LIMIT $2 OFFSET $3;

-- name: GetCategory :one
SELECT *
FROM category.category
WHERE id = $1 AND user_id = $2;

-- name: CreateCategory :one
INSERT INTO category.category (id, user_id, title, pinned, priority, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateCategory :one
UPDATE category.category
SET title      = $2,
    pinned     = $3,
    priority   = $4,
    updated_at = $5
WHERE id = $1 AND user_id = $6 RETURNING *;

-- name: DeleteCategory :exec
DELETE
FROM category.category
WHERE id = $1 AND user_id = $2;