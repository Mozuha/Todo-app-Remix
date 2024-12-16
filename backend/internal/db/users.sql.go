// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, user_id, username, email, password_hash, created_at, updated_at
`

type CreateUserParams struct {
	Email        string
	PasswordHash string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, user_id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, user_id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUserID = `-- name: GetUserByUserID :one
SELECT id, user_id, username, email, password_hash, created_at, updated_at FROM users WHERE user_id = $1
`

func (q *Queries) GetUserByUserID(ctx context.Context, userID pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUserID, userID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUsername = `-- name: UpdateUsername :exec
UPDATE users
SET username = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2
`

type UpdateUsernameParams struct {
	Username string
	ID       int32
}

func (q *Queries) UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error {
	_, err := q.db.Exec(ctx, updateUsername, arg.Username, arg.ID)
	return err
}