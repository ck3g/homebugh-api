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

	user = models.User{
		Email:             email,
		EncryptedPassword: encryptedPassword,
		CreatedAt:         time.Now(),
	}

	return user, nil
}
