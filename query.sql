-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;
