// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type AccessLog struct {
	ID           string
	UserID       sql.NullString
	BentoID      sql.NullString
	GroupID      sql.NullString
	BentoTokenID sql.NullString
	Action       string
	Details      interface{}
	AccessedAt   string
}

type Bento struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt string
	UpdatedAt string
}

type BentoIngridient struct {
	ID        string
	BentoID   string
	Name      string
	Value     []byte
	CreatedAt string
	UpdatedAt string
}

type BentoPermission struct {
	UserID    string
	BentoID   string
	Level     int64
	CreatedAt string
	UpdatedAt string
}

type BentoToken struct {
	ID         string
	BentoID    string
	TokenSalt  []byte
	CreatedBy  string
	CreatedAt  string
	LastUsedAt sql.NullString
	ExpiresAt  sql.NullString
}

type Group struct {
	ID        string
	Name      string
	OwnerID   string
	CreatedAt string
	UpdatedAt string
}

type GroupPermission struct {
	GroupID   string
	BentoID   string
	Level     int64
	CreatedAt string
	UpdatedAt string
}

type Jwt struct {
	ID        string
	UserID    string
	CreatedAt string
	ExpiresAt string
	TokenType string
}

type TotpRecoveryCode struct {
	ID        string
	UserID    string
	Code      string
	Used      bool
	CreatedAt string
}

type User struct {
	ID            string
	Email         string
	Password      string
	Nickname      string
	EmailVerified bool
	TotpSecret    sql.NullString
	TotpLocked    bool
	CreatedAt     string
	UpdatedAt     string
}

type UsersGroup struct {
	UserID    string
	GroupID   string
	CreatedAt string
}
