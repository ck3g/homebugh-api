package mysql

import (
	"errors"
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func TestUserConfirm(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	users := &UserModel{db}

	t.Run("Confirm non-confirmed user", func(t *testing.T) {
		id, err := users.Insert("user@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		_, err = users.Get(id)
		if err != nil {
			t.Fatal(err)
		}

		err = users.Confirm(id)
		if err != nil {
			t.Errorf("want error to be nil; got %s", err)
		}

		u, err := users.Get(id)
		if err != nil {
			t.Fatal(err)
		}

		if !u.ConfirmedAt.Valid {
			t.Errorf("want confirmed_at to be set; got nil")
		}
	})
}

func TestUserInsert(t *testing.T) {
	t.Run("successful insert", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		users := &UserModel{db}

		id, err := users.Insert("user@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		u, err := users.Get(id)
		if err != nil {
			t.Fatalf("expect to receive a user with ID %d", id)
		}

		if u.Email != "user@example.com" {
			t.Errorf("want email %s; got %s", "user@example.com", u.Email)
		}

		err = bcrypt.CompareHashAndPassword(u.EncryptedPassword, []byte("password"))
		if err != nil {
			t.Error("hashed password does not match")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		users := &UserModel{db}

		_, err := users.Insert("user@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}

		_, err = users.Insert("user@example.com", "password")
		if !errors.Is(err, models.ErrDuplicateEmail) {
			t.Errorf("want error %s; got %s", models.ErrDuplicateEmail, err)
		}
	})

}

func TestUserGet(t *testing.T) {
	t.Run("fetch existing user", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		users := &UserModel{db}

		id, err := users.Insert("user@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}
		defer users.Delete(id)

		u, err := users.Get(id)
		if err != nil {
			t.Fatalf("expect to receive user iwth ID %d", id)
		}

		if u.ID != id {
			t.Errorf("want ID %d; got %d", id, u.ID)
		}
	})

	t.Run("fetch non-existing user", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		users := &UserModel{db}

		_, err := users.Get(-1)
		if !errors.Is(err, models.ErrNoRecord) {
			t.Errorf("want error %s; got %s", models.ErrNoRecord, err)
		}
	})
}

func TestGetByEmail(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	users := &UserModel{db}
	id, err := users.Insert("user@example.com", "password")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name      string
		email     string
		wantID    int64
		wantError error
	}{
		{"Found with lower case email", "user@example.com", id, nil},
		{"Found with capital case email", "USER@EXAMPLE.COM", id, nil},
		{"Not found", "not@found.com", 0, models.ErrNoRecord},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := users.GetByEmail(tt.email)

			if u.ID != tt.wantID {
				t.Errorf("user not found: want ID %d; got ID: %d", tt.wantID, u.ID)
			}

			if err != tt.wantError {
				t.Errorf("invalid error returned: want %s; got %s", tt.wantError, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	users := &UserModel{db}

	id, err := users.Insert("user@example.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	u, err := users.Get(id)
	if err != nil {
		t.Errorf("want no errors; got %s", err)
	}

	err = users.Delete(u.ID)
	if err != nil {
		t.Errorf("want no errors: got %s", err)
	}

	u, err = users.Get(u.ID)
	if !errors.Is(err, models.ErrNoRecord) {
		t.Errorf("want error %s; found user with id %d", models.ErrNoRecord, u.ID)
	}
}

func TestAuthenticate(t *testing.T) {
	t.Run("successful auth", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		users := &UserModel{db}
		id, err := users.Insert("user@example.com", "password")
		if err != nil {
			panic(err)
		}

		users.Confirm(id)
		token, err := users.Authenticate("user@example.com", "password")
		if token == "" {
			t.Errorf("incorrect token: want token; got blank")
		}

		if err != nil {
			t.Errorf("wrong error returned: want nil; got %s", err)
		}

		// check if there is a new row in `auth_sessions` table with `token` in it
		sessions := &AuthSessionModel{db}
		s, err := sessions.GetByToken(token)
		if err != nil {
			t.Errorf("don't want errors; got %s", err)
		}

		if s.Token != token {
			t.Errorf("want session with token %s; got %s", token, s.Token)
		}

		if s.UserID != id {
			t.Errorf("want user_id %d; got %d", id, s.UserID)
		}
	})

	tests := []struct {
		name      string
		email     string
		password  string
		confirmed bool
		wantError error
	}{
		{"Wrong email", "no-user@example.com", "password", true, models.ErrNoRecord},
		{"Wrong password", "user@example.com", "wrong-pass", true, models.ErrWrongPassword},
		{"Not confirmed", "user@example.com", "password", false, models.ErrUserNotConfirmed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			users := &UserModel{db}
			id, err := users.Insert("user@example.com", "password")
			if err != nil {
				panic(err)
			}

			if tt.confirmed {
				if err := users.Confirm(id); err != nil {
					t.Fatal(err)
				}
			}

			token, err := users.Authenticate(tt.email, tt.password)

			if token != "" {
				t.Errorf("incorrect token: want blank; got %s", token)
			}

			if err != tt.wantError {
				t.Errorf("wrong error returned: want %v; got %v", tt.wantError, err)
			}

			// check there are no new rows in `auth_sessions` table
			sessions := &AuthSessionModel{db}
			_, err = sessions.GetByToken(token)
			if !errors.Is(err, models.ErrNoRecord) {
				t.Errorf("want error %s; got %s", models.ErrNoRecord, err)
			}
		})
	}
}
