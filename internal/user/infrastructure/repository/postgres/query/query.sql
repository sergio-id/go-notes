-- name: GetUser :one
SELECT *
FROM "user"."user"
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO "user"."user" (id, email, password, first_name, last_name, updated_at, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateUser :one
UPDATE "user"."user"
SET first_name = $2,
    last_name  = $3
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM "user"."user"
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM "user"."user"
WHERE email = $1;
