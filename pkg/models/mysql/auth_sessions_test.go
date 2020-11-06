package mysql

import (
	"errors"
	"testing"
	"time"

	"github.com/ck3g/homebugh-api/pkg/models"
)

func TestAuthSessionInsert(t *testing.T) {
	t.Run("successful insert", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		sessions := &AuthSessionModel{db}

		id, err := sessions.Insert(1, "unique-token")
		if err != nil {
			t.Fatal(err)
		}

		s, err := sessions.Get(id)
		if err != nil {
			t.Errorf("want to receive an auth session with ID %d", id)
		}

		if s.Token != "unique-token" {
			t.Errorf("want token %s; got %s", "unique-token", s.Token)
		}

		if s.UserID != 1 {
			t.Errorf("want user_id %d; got %d", 1, s.UserID)
		}

		days13 := time.Now().Add(time.Hour * 24 * 13).UTC()
		days15 := time.Now().Add(time.Hour * 24 * 15).UTC()
		within2Weeks := s.ExpiredAt.After(days13) && s.ExpiredAt.Before(days15)
		if !within2Weeks {
			t.Errorf("want expired_at to be in 2 weeks; got %s", s.ExpiredAt)
		}
	})

	t.Run("duplicated token", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		sessions := &AuthSessionModel{db}

		_, err := sessions.Insert(1, "unique-token")
		if err != nil {
			t.Fatal(err)
		}

		_, err = sessions.Insert(1, "unique-token")
		if !errors.Is(err, models.ErrDuplicateToken) {
			t.Errorf("want error %s; got %s", models.ErrDuplicateToken, err)
		}
	})
}

func TestAuthSessionGet(t *testing.T) {
	t.Run("fetch existing auth session", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		sessions := &AuthSessionModel{db}

		id, err := sessions.Insert(1, "token")
		if err != nil {
			t.Fatal(err)
		}

		s, err := sessions.Get(id)
		if err != nil {
			t.Errorf("expect to receive auth session with ID %d", id)
		}

		if s.ID != id {
			t.Errorf("want ID %d; got %d", id, s.ID)
		}
	})

	t.Run("fetch non-existing auth session", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		sessions := &AuthSessionModel{db}

		_, err := sessions.Get(-1)
		if !errors.Is(err, models.ErrNoRecord) {
			t.Errorf("want error %s; got %s", models.ErrNoRecord, err)
		}
	})

}

func TestAuthSessionDelete(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	sessions := &AuthSessionModel{DB: db}

	id, err := sessions.Insert(1, "token")
	if err != nil {
		t.Fatal(err)
	}

	s, err := sessions.Get(id)
	if err != nil {
		t.Errorf("want no errors; got %s", err)
	}

	err = sessions.Delete(s.ID)
	if err != nil {
		t.Errorf("want to errors: got %s", err)
	}

	s, err = sessions.Get(s.ID)
	if !errors.Is(err, models.ErrNoRecord) {
		t.Errorf("want error %s; found auth session with id %d", models.ErrNoRecord, s.ID)
	}
}
