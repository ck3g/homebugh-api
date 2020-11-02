package mysql

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserModel represents the MySQL data storage for users
type UserModel struct {
	DB *sql.DB
}

// Insert creates a new user in the database
func (m *UserModel) Insert(email, password string) (models.User, error) {
	user := models.User{}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	stmt := `INSERT INTO users (email, encrypted_password, password_salt, created_at)
	VALUES (?, ?, "", UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, email, string(encryptedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// Check if MySQL error is a email constraint violation
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return user, models.ErrDuplicateEmail
			}
		}

		return user, err
	}

	now := time.Now()
	user = models.User{
		Email:             email,
		EncryptedPassword: encryptedPassword,
		CreatedAt:         &now,
	}

	return user, nil
}

// Get fetches a user by ID. Returns an error if the user not found
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT id, email, encrypted_password, created_at, confirmed_at FROM users WHERE id = ?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Email, &u.EncryptedPassword, &u.CreatedAt, &u.ConfirmedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return u, nil
}
