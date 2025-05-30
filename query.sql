-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password, updated_at)
VALUES ($1, $2, NOW())
RETURNING *;

-- name: CreateEmailVerificationToken :one
INSERT INTO email_verification_tokens (email, verification_token, hashed_password_cache_key, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetEmailVerificationTokenByEmail :one
SELECT * FROM email_verification_tokens
WHERE email = $1;

-- name: DeleteEmailVerificationTokenByID :exec
DELETE FROM email_verification_tokens
WHERE id = $1;
