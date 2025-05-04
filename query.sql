-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
