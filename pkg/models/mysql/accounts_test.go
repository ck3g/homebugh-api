package mysql

import (
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"
)

func TestAccountInsert(t *testing.T) {
	t.Run("successful insert", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		accounts := &AccountModel{db}
		id, err := accounts.Insert("Bank", 1, 1, "active", true)
		if err != nil {
			t.Fatal(err)
		}

		filters := models.Filters{
			Page:     1,
			PageSize: 20,
		}

		all, _, err := accounts.All(1, filters)
		if err != nil {
			t.Fatal(err)
		}

		if len(all) != 1 {
			t.Errorf("want 1 account; got %d", len(all))
			return
		}

		account := all[0]
		if account.ID != id {
			t.Errorf("want ID %d; got %d", id, account.ID)
		}

		if account.Name != "Bank" {
			t.Errorf("want Name %s; got %s", "Bank", account.Name)
		}

		if account.CurrencyID != 1 {
			t.Errorf("want CurrencyID %d; got %d", 1, account.CurrencyID)
		}

		if account.Status != "active" {
			t.Errorf("want status %s; got %s", "active", account.Status)
		}

		if !account.ShowInSummary {
			t.Errorf("want show in summary %t; got %t", true, account.ShowInSummary)
		}
	})
}
