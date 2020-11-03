package mysql

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

const dsn = "root@/homebugh_test?parseTime=true"

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func TestInsert(t *testing.T) {
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	users := &UserModel{DB: db}

	t.Run("successful insert", func(t *testing.T) {
		id, err := users.Insert("user@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}
		defer users.Delete(id)

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
		id, err := users.Insert("user@example.com", "password")
		if err != nil {
			t.Fatal(err)
		}
		defer users.Delete(id)

		_, err = users.Insert("user@example.com", "password")
		if !errors.Is(err, models.ErrDuplicateEmail) {
			t.Errorf("want error %s; got %s", models.ErrDuplicateEmail, err)
		}
	})

}

func TestGet(t *testing.T) {
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	users := &UserModel{DB: db}

	t.Run("fetch existing user", func(t *testing.T) {
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
		_, err := users.Get(-1)
		if !errors.Is(err, models.ErrNoRecord) {
			t.Errorf("want error %s; got %s", models.ErrNoRecord, err)
		}
	})
}

func TestDelete(t *testing.T) {
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	users := &UserModel{DB: db}

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
