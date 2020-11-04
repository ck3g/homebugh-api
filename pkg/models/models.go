package models

import (
	"errors"
	"time"
)

var (
	// ErrDuplicateEmail received when user with the same email already exists
	ErrDuplicateEmail = errors.New("models: duplicate email")

	// ErrNoRecord received when the record not found in the database
	ErrNoRecord = errors.New("models: record not found")

	// ErrWrongPassword recived when the password is incorrect
	ErrWrongPassword = errors.New("users: wrong password")
)

// User represents a user data
type User struct {
	ID                int64
	Email             string
	EncryptedPassword []byte
	CreatedAt         *time.Time
	ConfirmedAt       *time.Time
}

// UserStorage defines interface to storing and retrieving user data
type UserStorage interface {
	Insert(email, password string) (int64, error)
	Get(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	Delete(id int64) error
	Authenticate(email, password string) (string, error)
}
