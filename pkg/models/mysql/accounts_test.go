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

func TestAccountAll(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	accounts := &AccountModel{db}

	accountsData := map[int]struct {
		name          string
		userID        int64
		currencyID    int64
		status        string
		showInSummary bool
	}{
		0: {"Bank", 1, 1, "active", true},
		1: {"Bank", 2, 1, "active", true},
		2: {"Cash", 1, 1, "active", true},
	}
	ids := map[int]int64{}

	for i, a := range accountsData {
		id, err := accounts.Insert(a.name, a.userID, a.currencyID, a.status, a.showInSummary)
		if err != nil {
			t.Fatal(err)
		}

		ids[i] = id
	}

	tests := []struct {
		name             string
		userID           int64
		filters          models.Filters
		wantCount        int
		wantTotalRecords int
		wantAccounts     []*models.Account
	}{
		{
			name:             "successful fetch for specific user",
			userID:           1,
			filters:          models.Filters{Page: 1, PageSize: 20},
			wantCount:        2,
			wantTotalRecords: 2,
			wantAccounts: []*models.Account{
				{ID: ids[0], Name: "Bank", CurrencyID: 1, Status: "active", ShowInSummary: true},
				{ID: ids[2], Name: "Cash", CurrencyID: 1, Status: "active", ShowInSummary: true},
			},
		},
		{
			name:             "accounts from the first page",
			userID:           1,
			filters:          models.Filters{Page: 1, PageSize: 1},
			wantCount:        1,
			wantTotalRecords: 2,
			wantAccounts: []*models.Account{
				{ID: ids[0], Name: "Bank", CurrencyID: 1, Status: "active", ShowInSummary: true},
			},
		},
		{
			name:             "accounts from the second page",
			userID:           1,
			filters:          models.Filters{Page: 2, PageSize: 1},
			wantCount:        1,
			wantTotalRecords: 2,
			wantAccounts: []*models.Account{
				{ID: ids[2], Name: "Cash", CurrencyID: 1, Status: "active", ShowInSummary: true},
			},
		},
		{
			name:             "accounts from the greater than last page",
			userID:           1,
			filters:          models.Filters{Page: 100, PageSize: 1},
			wantCount:        0,
			wantTotalRecords: 2,
			wantAccounts:     []*models.Account{},
		},
		{
			name:             "accounts from the negative page",
			userID:           1,
			filters:          models.Filters{Page: -1, PageSize: 1},
			wantCount:        1,
			wantTotalRecords: 2,
			wantAccounts: []*models.Account{
				{ID: ids[0], Name: "Bank", CurrencyID: 1, Status: "active", ShowInSummary: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all, metadata, err := accounts.All(tt.userID, tt.filters)
			if err != nil {
				t.Fatal(err)
			}

			if len(all) != tt.wantCount {
				t.Errorf("want %d accounts; got %d", tt.wantCount, len(all))
			}

			if metadata.TotalRecords != tt.wantTotalRecords {
				t.Errorf("want %d total records; got %d", tt.wantTotalRecords, metadata.TotalRecords)
			}

			for i, got := range all {
				wantAccount := tt.wantAccounts[i]

				if got.ID != wantAccount.ID {
					t.Errorf("want ID %d; got %d", wantAccount.ID, got.ID)
				}

				if got.Name != wantAccount.Name {
					t.Errorf("want Name %s; got %s", wantAccount.Name, got.Name)
				}

				if got.CurrencyID != wantAccount.CurrencyID {
					t.Errorf("want CurrencyID %d; got %d", wantAccount.CurrencyID, got.CurrencyID)
				}

				if got.Status != wantAccount.Status {
					t.Errorf("want Status %s; got %s", wantAccount.Status, got.Status)
				}

				if got.ShowInSummary != wantAccount.ShowInSummary {
					t.Errorf("want ShowInSummary %t; got %t", wantAccount.ShowInSummary, got.ShowInSummary)
				}
			}
		})
	}
}
