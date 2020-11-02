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
)

// User represents a user data
type User struct {
	ID                int
	Email             string
	EncryptedPassword []byte
	CreatedAt         *time.Time
	ConfirmedAt       *time.Time
}

// UserStorage defines interface to storing and retrieving user data
type UserStorage interface {
	Insert(email, password string) (User, error)
	Get(id int) (*User, error)
}
