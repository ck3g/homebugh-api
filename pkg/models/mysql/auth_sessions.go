package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/go-sql-driver/mysql"
)

// AuthSessionModel represents the MySQL data storage for authentication sessions
type AuthSessionModel struct {
	DB *sql.DB
}

// Insert create a new authentication session
func (m *AuthSessionModel) Insert(userID int64, token string) (int64, error) {
	var id int64

	stmt := `INSERT INTO auth_sessions (user_id, token, expired_at, created_at, updated_at)
	VALUES(?, ?, DATE_ADD(UTC_TIMESTAMP(), INTERVAL 2 WEEK), UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	res, err := m.DB.Exec(stmt, userID, token)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			// Check if MySQL error is a token constraint violation
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "auth_sessions_on_token") {
				return id, models.ErrDuplicateToken
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

// Get retrieves auth_session record by its ID
func (m *AuthSessionModel) Get(id int64) (*models.AuthSession, error) {
	s := &models.AuthSession{}

	stmt := `SELECT id, user_id, token, created_at, expired_at FROM auth_sessions WHERE id = ?`

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.UserID, &s.Token, &s.CreatedAt, &s.ExpiredAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}

		return nil, err
	}

	return s, nil
}

// Delete removed authentication session by its ID
func (m *AuthSessionModel) Delete(id int64) error {
	stmt := `DELETE FROM auth_sessions WHERE id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
