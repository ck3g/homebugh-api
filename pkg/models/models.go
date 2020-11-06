package models

import (
	"errors"
	"time"
)

var (
	// ErrDuplicateEmail received when user with the same email already exists
	ErrDuplicateEmail = errors.New("models: duplicate email")

	// ErrDuplicateToken received when auth_session with the same token already exists
	ErrDuplicateToken = errors.New("models: duplicate token")

	// ErrNoRecord received when the record not found in the database
	ErrNoRecord = errors.New("models: record not found")

	// ErrWrongPassword recived when the password is incorrect
	ErrWrongPassword = errors.New("users: wrong password")
)

// AuthSession represents authentication session data
type AuthSession struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiredAt *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

// AuthSessionStorage contains information about current API authentications
type AuthSessionStorage interface {
	Insert(userID int64, token string) (int64, error)
	Get(id int64) (*AuthSession, error)
	GetByToken(token string) (*AuthSession, error)
	Delete(id int64) error
}

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
