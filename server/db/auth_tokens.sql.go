// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: auth_tokens.sql

package db

import (
	"context"
)

const deletAuthTokenById = `-- name: DeletAuthTokenById :exec
DELETE FROM auth_tokens WHERE id = ?
`

func (q *Queries) DeletAuthTokenById(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deletAuthTokenById, id)
	return err
}

const deleteAllTokensByTypeAndUserID = `-- name: DeleteAllTokensByTypeAndUserID :exec
DELETE FROM auth_tokens WHERE user_id = ? AND token_type = ?
`

type DeleteAllTokensByTypeAndUserIDParams struct {
	UserID    string `db:"user_id" json:"user_id"`
	TokenType string `db:"token_type" json:"token_type"`
}

func (q *Queries) DeleteAllTokensByTypeAndUserID(ctx context.Context, arg DeleteAllTokensByTypeAndUserIDParams) error {
	_, err := q.db.ExecContext(ctx, deleteAllTokensByTypeAndUserID, arg.UserID, arg.TokenType)
	return err
}

const deleteUserAuthTokens = `-- name: DeleteUserAuthTokens :exec
DELETE FROM auth_tokens WHERE user_id = ?
`

func (q *Queries) DeleteUserAuthTokens(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, deleteUserAuthTokens, userID)
	return err
}

const existsAuthTokenById = `-- name: ExistsAuthTokenById :one
SELECT EXISTS(SELECT 1 FROM auth_tokens WHERE id = ?)
`

func (q *Queries) ExistsAuthTokenById(ctx context.Context, id string) (int64, error) {
	row := q.db.QueryRowContext(ctx, existsAuthTokenById, id)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getAuthTokenById = `-- name: GetAuthTokenById :one
SELECT id, user_id, created_at, expires_at, token_type FROM auth_tokens
WHERE id = ?
`

func (q *Queries) GetAuthTokenById(ctx context.Context, id string) (AuthToken, error) {
	row := q.db.QueryRowContext(ctx, getAuthTokenById, id)
	var i AuthToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.TokenType,
	)
	return i, err
}

const getUserAuthTokens = `-- name: GetUserAuthTokens :many
SELECT id, user_id, created_at, expires_at, token_type FROM auth_tokens
WHERE user_id = ?
`

func (q *Queries) GetUserAuthTokens(ctx context.Context, userID string) ([]AuthToken, error) {
	rows, err := q.db.QueryContext(ctx, getUserAuthTokens, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuthToken
	for rows.Next() {
		var i AuthToken
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.CreatedAt,
			&i.ExpiresAt,
			&i.TokenType,
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

const newAuthToken = `-- name: NewAuthToken :one
INSERT INTO auth_tokens
(user_id, created_at, expires_at, token_type)
VALUES
(?, ?, ?, ?)
RETURNING id, user_id, created_at, expires_at, token_type
`

type NewAuthTokenParams struct {
	UserID    string `db:"user_id" json:"user_id"`
	CreatedAt string `db:"created_at" json:"created_at"`
	ExpiresAt string `db:"expires_at" json:"expires_at"`
	TokenType string `db:"token_type" json:"token_type"`
}

func (q *Queries) NewAuthToken(ctx context.Context, arg NewAuthTokenParams) (AuthToken, error) {
	row := q.db.QueryRowContext(ctx, newAuthToken,
		arg.UserID,
		arg.CreatedAt,
		arg.ExpiresAt,
		arg.TokenType,
	)
	var i AuthToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CreatedAt,
		&i.ExpiresAt,
		&i.TokenType,
	)
	return i, err
}
