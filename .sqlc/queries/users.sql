-- name: CreateUser :one
INSERT INTO users
(email, password, nickname, created_at, updated_at)
VALUES
(?, ?, ?, ?, ?)
RETURNING id;

-- name: ExistsUserWithEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = ?);

-- name: DeleteUserById :exec
DELETE FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT
    id,
    email,
    email_verified,
    password,
    nickname,
    totp_secret,
    created_at,
    updated_at
FROM users
WHERE email = ?;

-- name: SetUserEmailVerifiedStatus :exec
UPDATE users SET email_verified = ? WHERE id = ?;
