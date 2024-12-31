// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users
(email, password, nickname, token_salt, created_at, updated_at)
VALUES
(?, ?, ?, ?, ?, ?)
RETURNING id, email_verified
`

type CreateUserParams struct {
	Email     string
	Password  string
	Nickname  string
	TokenSalt []byte
	CreatedAt string
	UpdatedAt string
}

type CreateUserRow struct {
	ID            string
	EmailVerified bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.Password,
		arg.Nickname,
		arg.TokenSalt,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i CreateUserRow
	err := row.Scan(&i.ID, &i.EmailVerified)
	return i, err
}

const existsUserWithEmail = `-- name: ExistsUserWithEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)
`

func (q *Queries) ExistsUserWithEmail(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, existsUserWithEmail, email)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}