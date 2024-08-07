package store

import (
	"database/sql"
	"os"
	"time"
)

// User represents a real user store in the database.
type User struct {
	Id            string
	Email         string
	Password      string
	Name          string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Delete removes a user from the database using the id.
//
// You must call tx.Commit for it to take effect.
func (u *User) Delete(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec("DELETE FROM users WHERE id = $1;", u.Id)
}

// DEPRACATED: this method will try to update the columns in users
// and rehash the password while doing so, which causes issues since the
// password is not the same anymore. This leads to users unable to access
// their accounts unless they reset their password.
//
// Update updates the user in the database with the current field values of the struct.
//
// This does not update the Id, CreatedAt, and UpdatedAt fields.
//
// You must call tx.Commit for it to take effect.
func (u *User) Update(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec("UPDATE users SET email = $1, password = crypt($2, gen_salt($6)), name = $3, email_verified = $4 WHERE id = $5;", u.Email, u.Password, u.Name, u.EmailVerified, u.Id, os.Getenv("PASS_ENCRYPT_ALGO"))
}

// Sets the value for the email_verified column in the users table.
// Will set the User.EmailVerified field to the boolean value 'v'.
//
// You must call tx.Commit() for it to take effect.
func (u *User) SetEmailVerifiedTx(tx *sql.Tx, v bool) (sql.Result, error) {
	res, err := tx.Exec("UPDATE users SET email_verified = $1 WHERE id = $2;", v, u.Id)
	if err != nil {
		return nil, err
	}
	u.EmailVerified = v
	return res, nil
}

// NewUser creates a new user with the given information.
// This function will save the new user in the database.
func NewUser(email, password, name string) (*User, error) {
	user := User{
		Email:    email,
		Password: password,
		Name:     name,
	}
	row := db.QueryRow("INSERT INTO users (email, password, name) VALUES ($1, crypt($2, gen_salt($3)), $4) RETURNING id, email_verified, created_at, updated_at;", email, password, os.Getenv("PASS_ENCRYPT_ALGO"), name)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	err = row.Scan(
		&user.Id,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsUserWithEmail checks if there is a user with the given email stored in the database.
func ExistsUserWithEmail(email string) (bool, error) {
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email)
	err := row.Err()
	if err != nil {
		return false, err
	}
	var exists bool
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetUserWithId retrieves a user with the given id.
// An error is returned if nothing is found.
func GetUserWithId(id string) (*User, error) {
	row := db.QueryRow("SELECT id, email, password, name, email_verified, created_at, updated_at FROM users WHERE id = $1;", id)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := User{}
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserWithEmail retrieves a user in the database with the given email.
func GetUserWithEmail(email string) (*User, error) {
	row := db.QueryRow("SELECT id, email, password, name, email_verified, created_at, updated_at FROM users WHERE email = $1;", email)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := User{}
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserWithEmailAndPassword tries to match a user with the given email and password.
// Usage for signining a user. Its better to use psql to match the password in the query.
func GetUserWithEmailAndPassword(email, password string) (*User, error) {
	row := db.QueryRow(
		"SELECT id, email, password, name, email_verified, created_at, updated_at FROM users WHERE email = $1 AND password = crypt($2, password);",
		email,
		password,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := User{}
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
