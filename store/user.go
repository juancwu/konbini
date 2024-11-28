package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string
	Email         string
	Nickname      *string
	PasswordHash  string
	IsActive      bool
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastLoginAt   string
}

const get_USER_BY_ID_SQL = `SELECT
    id,
    email,
    nickname,
    password_hash,
    is_active,
    email_verified,
    created_at,
    updated_at,
    last_login_at
FROM users WHERE id = $1;
`

// Gets a user by ID.
func GetUserByID(ctx context.Context, db *sql.DB, id string) (*User, error) {
	row := db.QueryRowContext(ctx, get_USER_BY_ID_SQL, id)
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.PasswordHash,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

const new_USER_SQL = `INSERT INTO users (id, email, password_hash, nickname) VALUES (?, ?, ?, ?);`

// Creates a new user with given email and password.
// This method will hash the password so DO NOT hash
// the password when calling the function.
func NewUser(ctx context.Context, db *sql.DB, email, password string, nickname *string) (string, error) {
	// generate uuid
	id := uuid.NewString()
	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	_, err = db.ExecContext(ctx, new_USER_SQL, id, email, string(passwordHash), nickname)
	if err != nil {
		return "", err
	}

	return id, nil
}

const exists_USER_SQL = `SELECT COUNT(*) FROM users WHERE email = ?;`

// ExistsUserWithEmail checks if a user with the given email exists or not.
func ExistsUserWithEmail(ctx context.Context, db *sql.DB, email string) (bool, error) {
	var count int
	row := db.QueryRowContext(ctx, exists_USER_SQL, email)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
