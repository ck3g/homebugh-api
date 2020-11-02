package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserModel represents the MySQL data storage for users
type UserModel struct {
	DB *sql.DB
}

// Insert creates a new user in the database
func (m *UserModel) Insert(email, password string) (int64, error) {
	var id int64

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return id, err
	}

	stmt := `INSERT INTO users (email, encrypted_password, password_salt, created_at)
	VALUES (?, ?, "", UTC_TIMESTAMP())`

	res, err := m.DB.Exec(stmt, email, string(encryptedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// Check if MySQL error is a email constraint violation
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return id, models.ErrDuplicateEmail
			}
		}

		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

// Get fetches a user by ID. Returns an error if the user not found
func (m *UserModel) Get(id int64) (*models.User, error) {
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

// Delete removes a user by ID
func (m *UserModel) Delete(id int64) error {
	stmt := `DELETE FROM users WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
