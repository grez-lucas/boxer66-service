// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createEmailVerificationToken = `-- name: CreateEmailVerificationToken :one
INSERT INTO email_verification_tokens (email, verification_token, hashed_password_cache_key, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, email, verification_token, hashed_password_cache_key, token_type, created_at, expires_at
`

type CreateEmailVerificationTokenParams struct {
	Email                  string    `json:"email"`
	VerificationToken      string    `json:"verification_token"`
	HashedPasswordCacheKey string    `json:"hashed_password_cache_key"`
	ExpiresAt              time.Time `json:"expires_at"`
}

func (q *Queries) CreateEmailVerificationToken(ctx context.Context, arg CreateEmailVerificationTokenParams) (EmailVerificationToken, error) {
	row := q.db.QueryRow(ctx, createEmailVerificationToken,
		arg.Email,
		arg.VerificationToken,
		arg.HashedPasswordCacheKey,
		arg.ExpiresAt,
	)
	var i EmailVerificationToken
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.VerificationToken,
		&i.HashedPasswordCacheKey,
		&i.TokenType,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password, updated_at)
VALUES ($1, $2, NOW())
RETURNING id, email, password, created_at, updated_at
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteEmailVerificationTokenByID = `-- name: DeleteEmailVerificationTokenByID :exec
DELETE FROM email_verification_tokens
WHERE id = $1
`

func (q *Queries) DeleteEmailVerificationTokenByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteEmailVerificationTokenByID, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, email, password, created_at, updated_at FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEmailVerificationTokenByEmail = `-- name: GetEmailVerificationTokenByEmail :one
SELECT id, email, verification_token, hashed_password_cache_key, token_type, created_at, expires_at FROM email_verification_tokens
WHERE email = $1
`

func (q *Queries) GetEmailVerificationTokenByEmail(ctx context.Context, email string) (EmailVerificationToken, error) {
	row := q.db.QueryRow(ctx, getEmailVerificationTokenByEmail, email)
	var i EmailVerificationToken
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.VerificationToken,
		&i.HashedPasswordCacheKey,
		&i.TokenType,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, created_at, updated_at FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, password, created_at, updated_at FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
