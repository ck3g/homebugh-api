package mock

import (
	"time"

	"github.com/ck3g/homebugh-api/pkg/models"
)

var (
	twoWeeksFromNow = time.Now().Add(time.Hour * 24 * 14)
	twoDaysAgo      = time.Now().Add(-time.Hour * 48)
	now             = time.Now()

	mockAuthSession = &models.AuthSession{
		ID:        1,
		UserID:    1,
		Token:     "valid-token",
		ExpiredAt: &twoWeeksFromNow,
		CreatedAt: &twoDaysAgo,
		UpdatedAt: &now,
	}
)

// AuthSessionModel represents mocked AuthSessionModel
type AuthSessionModel struct{}

// Insert mocks auth session insert method
func (m *AuthSessionModel) Insert(userID int64, token string) (int64, error) {
	return mockAuthSession.ID, nil
}

// Get fetches mock auth session by ID
func (m *AuthSessionModel) Get(id int64) (*models.AuthSession, error) {
	if id == mockAuthSession.ID {
		return mockAuthSession, nil
	}

	return nil, models.ErrNoRecord
}

// GetByToken fetches mock auth session by token
func (m *AuthSessionModel) GetByToken(token string) (*models.AuthSession, error) {
	if token == mockAuthSession.Token {
		return mockAuthSession, nil
	}

	return nil, models.ErrNoRecord
}

// Delete drops mocked auth session
func (m *AuthSessionModel) Delete(id int64) error {
	return nil
}
