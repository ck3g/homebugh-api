package mock

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/ck3g/homebugh-api/pkg/models"
	"github.com/go-sql-driver/mysql"
)

var (
	oneDayAgo = time.Now().Add(-time.Hour * 24)
	nullTime  = mysql.NullTime{
		Time:  oneDayAgo,
		Valid: false,
	}
	nonNullTime = mysql.NullTime{
		Time:  oneDayAgo,
		Valid: true,
	}
	encryptedPassword, _ = bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockConfirmedUser = &models.User{
		ID:                1,
		Email:             "user@example.com",
		EncryptedPassword: encryptedPassword,
		CreatedAt:         &oneDayAgo,
		UpdatedAt:         &oneDayAgo,
		ConfirmedAt:       nonNullTime,
	}

	mockNotConfirmedUser = &models.User{
		ID:                2,
		Email:             "not-confirmed@example.com",
		EncryptedPassword: encryptedPassword,
		CreatedAt:         &oneDayAgo,
		UpdatedAt:         &oneDayAgo,
		ConfirmedAt:       nullTime,
	}
)

// UserModel represents mock UserModel
type UserModel struct{}

// Authenticate mock authenticate user
func (m *UserModel) Authenticate(email, password string) (string, error) {
	if email == mockNotConfirmedUser.Email {
		return "", models.ErrUserNotConfirmed
	}

	if email != mockConfirmedUser.Email {
		return "", models.ErrNoRecord
	}

	err := bcrypt.CompareHashAndPassword(mockConfirmedUser.EncryptedPassword, []byte(password))
	if err != nil {
		return "", models.ErrWrongPassword
	}

	return "valid-token", nil
}

// Confirm mock user confirmation
func (m *UserModel) Confirm(id int64) error {
	return nil
}

// Get mock fetching user by ID
func (m *UserModel) Get(id int64) (*models.User, error) {
	if id == mockConfirmedUser.ID {
		return mockConfirmedUser, nil
	}

	if id == mockNotConfirmedUser.ID {
		return mockNotConfirmedUser, nil
	}

	return nil, models.ErrNoRecord
}

// GetByEmail mock fetching user by Email
func (m *UserModel) GetByEmail(email string) (*models.User, error) {
	if email == mockConfirmedUser.Email {
		return mockConfirmedUser, nil
	}

	if email == mockNotConfirmedUser.Email {
		return mockNotConfirmedUser, nil
	}

	return nil, models.ErrNoRecord
}

// Insert mock insert a new user
func (m *UserModel) Insert(email, password string) (int64, error) {
	if email == mockConfirmedUser.Email {
		return mockConfirmedUser.ID, nil
	}

	if email == mockNotConfirmedUser.Email {
		return mockNotConfirmedUser.ID, nil
	}

	return 0, models.ErrDuplicateEmail
}

// Delete mock delete user
func (m *UserModel) Delete(id int64) error {
	return nil
}
