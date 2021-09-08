package models

import (
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
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

	// ErrUserNotConfirmed received when user is not confirmed
	ErrUserNotConfirmed = errors.New("users: not confirmed")
)

type Models struct {
	Accounts     AccountStorage
	AuthSessions AuthSessionStorage
	Categories   CategoryStorage
	Users        UserStorage
}

// Account represents a user's accounts
type Account struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Balance       float64 `json:"balance"`
	CurrencyID    int64   `json:"currency_id"` // TODO: change to currency
	Status        string  `json:"status"`
	ShowInSummary bool    `json:"show_in_summary"`
}

// AccountStorate contains information about user's accounts
type AccountStorage interface {
	Insert(name string, userID int64, currencyID int64, status string, showInSummary bool) (int64, error)
	All(userID int64, filters Filters) ([]*Account, Metadata, error)
}

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

// CategoryType represends a category type: with ID = 1 - income, ID = 2 - expense
type CategoryType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Category represents an expense category
type Category struct {
	ID           int64        `json:"id"`
	Name         string       `json:"name"`
	CategoryType CategoryType `json:"category_type"`
	UserID       int64        `json:"-"`
	Inactive     bool         `json:"inactive"`
	UpdatedAt    *time.Time   `json:"-"`
}

type CategoryStorage interface {
	Insert(name string, categoryType CategoryType, userID int64, inactive bool) (int64, error)
	All(userID int64, filters Filters) ([]*Category, Metadata, error)
}

// User represents a user data
type User struct {
	ID                int64
	Email             string
	EncryptedPassword []byte
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	ConfirmedAt       mysql.NullTime // https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
}

// UserStorage defines interface to storing and retrieving user data
type UserStorage interface {
	Authenticate(email, password string) (string, error)
	Confirm(id int64) error
	Get(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	Insert(email, password string) (int64, error)
	Delete(id int64) error
}
