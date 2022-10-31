package mysql

import (
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"
)

func TestTransaction_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	t.Run("successful insert", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		categories := &CategoryModel{db}
		income := models.CategoryType{ID: 1, Name: "income"}
		categoryID, err := categories.Insert("Salary", income, 1, false)
		if err != nil {
			t.Fatal(err)
		}

		accounts := &AccountModel{db}
		accountID, err := accounts.Insert("Cash", 1, 1, "active", true)
		if err != nil {
			t.Fatal(err)
		}

		transactions := &TransactionModel{db}
		id, err := transactions.Insert(100, "some income", 1, categoryID, accountID)
		if err != nil {
			t.Fatal(err)
		}

		filters := models.Filters{
			Page:     1,
			PageSize: 20,
		}

		all, _, err := transactions.All(1, filters)
		if err != nil {
			t.Fatal(err)
		}

		if len(all) != 1 {
			t.Errorf("want 1 transaction; got %d", len(all))
		}

		transaction := all[0]
		if transaction.ID != id {
			t.Errorf("want ID %d; got %d", id, transaction.ID)
		}

		if transaction.Amount != 100.0 {
			t.Errorf("want Amount %f; got %f", 100.0, transaction.Amount)
		}

		if transaction.Comment != "some income" {
			t.Errorf("want Comment %s; got %s", "some income", transaction.Comment)
		}

		if transaction.CreatedAt == nil {
			t.Errorf("want CreatedAt to be present; got nil")
		}

		if transaction.Category.ID != categoryID {
			t.Errorf("want Category.ID %d; got %d", categoryID, transaction.Category.ID)
		}

		if transaction.Category.Name != "Salary" {
			t.Errorf("want Category.Name %s; got %s", "Salary", transaction.Category.Name)
		}

		if transaction.Category.Inactive != false {
			t.Errorf("want Category.Inactive %t; got %t", false, transaction.Category.Inactive)
		}

		if transaction.Category.CategoryType.ID != 1 {
			t.Errorf("want CategoryType.ID %d; got %d", 1, transaction.Category.CategoryType.ID)
		}

		if transaction.Category.CategoryType.Name != "income" {
			t.Errorf("want CategoryType.Name %s; got %s", "income", transaction.Category.CategoryType.Name)
		}

		if transaction.Account.ID != accountID {
			t.Errorf("want Account.ID %d; got %d", accountID, transaction.Account.ID)
		}

		if transaction.Account.Name != "Cash" {
			t.Errorf("want Account.Name %s; got %s", "Cash", transaction.Account.Name)
		}

		if transaction.Account.Status != "active" {
			t.Errorf("want Account.Status %s; got %s", "active", transaction.Account.Status)
		}

		if transaction.Account.ShowInSummary != true {
			t.Errorf("want Account.ShowInSummary %t; got %t", true, transaction.Account.ShowInSummary)
		}

		if transaction.Account.Currency.ID != 1 {
			t.Errorf("want Currency.ID %d; got %d", 1, transaction.Account.Currency.ID)
		}

		if transaction.Account.Currency.Name != "Euro" {
			t.Errorf("want Currency.Name %s; got %s", "Euro", transaction.Account.Currency.Name)
		}

		if transaction.Account.Currency.Unit != "€" {
			t.Errorf("want Currency.Unit %s; got %s", "€", transaction.Account.Currency.Unit)
		}
	})
}

func TestTransaction_All(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	db, teardown := newTestDB(t)
	defer teardown()

	currency := models.Currency{ID: 1, Name: "Euro", Unit: "€"}

	// Insert accounts
	accounts := &AccountModel{db}
	allAccounts := map[int]*models.Account{
		0: {Name: "Bank", UserID: 1, Currency: currency, Status: "active", ShowInSummary: true},
		1: {Name: "Bank", UserID: 2, Currency: currency, Status: "active", ShowInSummary: true},
		2: {Name: "Cash", UserID: 1, Currency: currency, Status: "active", ShowInSummary: true},
	}
	accountIDs := map[int]int64{}

	for i, a := range allAccounts {
		id, err := accounts.Insert(a.Name, a.UserID, a.Currency.ID, a.Status, a.ShowInSummary)
		if err != nil {
			t.Fatal(err)
		}

		accountIDs[i] = id
		allAccounts[i].ID = id
	}

	// Insert categories
	categories := &CategoryModel{db}
	income := models.CategoryType{ID: 1, Name: "income"}
	expense := models.CategoryType{ID: 2, Name: "expense"}

	allCategories := map[int]*models.Category{
		0: {Name: "Food", CategoryType: expense, UserID: 2, Inactive: true},
		1: {Name: "Clothes", CategoryType: expense, UserID: 1},
		2: {Name: "Salary", CategoryType: income, UserID: 1},
	}
	categoryIDs := map[int]int64{}

	for i, c := range allCategories {
		id, err := categories.Insert(c.Name, c.CategoryType, c.UserID, c.Inactive)
		if err != nil {
			t.Fatal(err)
		}

		categoryIDs[i] = id
		allCategories[i].ID = id
	}

	// Insert transactions
	transactions := &TransactionModel{db}

	allTransactions := map[int]*models.Transaction{
		0: {Amount: 10, Comment: "comment1", UserID: 1, Category: *allCategories[1], Account: *allAccounts[0]},
		1: {Amount: 11, Comment: "comment2", UserID: 2, Category: *allCategories[0], Account: *allAccounts[1]},
		2: {Amount: 12, Comment: "comment3", UserID: 1, Category: *allCategories[2], Account: *allAccounts[2]},
	}

	ids := map[int]int64{}

	for i, tr := range allTransactions {
		id, err := transactions.Insert(tr.Amount, tr.Comment, tr.UserID, tr.Category.ID, tr.Account.ID)
		if err != nil {
			t.Fatal(err)
		}

		ids[i] = id
	}

	tests := []struct {
		name             string
		wantTransactions []*models.Transaction
		wantCount        int
		wantTotalRecords int
		userID           int64
		filters          models.Filters
	}{
		{
			name: "successful fetch for a specific user",
			wantTransactions: []*models.Transaction{
				{ID: ids[2], Amount: 12, Comment: "comment3"},
				{ID: ids[0], Amount: 10, Comment: "comment1"},
			},
			wantCount:        2,
			wantTotalRecords: 2,
			userID:           1,
			filters:          models.Filters{Page: 1, PageSize: 20},
		},
		{
			name: "transactions from the first page",
			wantTransactions: []*models.Transaction{
				{ID: ids[2], Amount: 12, Comment: "comment3"},
			},
			wantCount:        1,
			wantTotalRecords: 2,
			userID:           1,
			filters:          models.Filters{Page: 1, PageSize: 1},
		},
		{
			name: "transactions from the second page",
			wantTransactions: []*models.Transaction{
				{ID: ids[0], Amount: 10, Comment: "comment1"},
			},
			wantCount:        1,
			wantTotalRecords: 2,
			userID:           1,
			filters:          models.Filters{Page: 2, PageSize: 1},
		},
		{
			name:             "transactions from the greater than last page",
			wantTransactions: []*models.Transaction{},
			wantCount:        0,
			wantTotalRecords: 2,
			userID:           1,
			filters:          models.Filters{Page: 100, PageSize: 1},
		},
		{
			name: "transactions from the negative page",
			wantTransactions: []*models.Transaction{
				{ID: ids[2], Amount: 12, Comment: "comment3"},
			},
			wantCount:        1,
			wantTotalRecords: 2,
			userID:           1,
			filters:          models.Filters{Page: -1, PageSize: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all, metadata, err := transactions.All(tt.userID, tt.filters)
			if err != nil {
				t.Fatal(err)
			}

			if len(all) != tt.wantCount {
				t.Errorf("want %d transactions; got %d", tt.wantCount, len(all))
			}

			if metadata.TotalRecords != tt.wantTotalRecords {
				t.Errorf("want %d total records; got %d", tt.wantTotalRecords, metadata.TotalRecords)
			}

			for i, got := range all {
				if got.ID != tt.wantTransactions[i].ID {
					t.Errorf("want ID %d; got %d", tt.wantTransactions[i].ID, got.ID)
				}

				if got.Amount != tt.wantTransactions[i].Amount {
					t.Errorf("want Amount %f; got %f", tt.wantTransactions[i].Amount, got.Amount)
				}

				if got.Comment != tt.wantTransactions[i].Comment {
					t.Errorf("want Comment %s; got %s", tt.wantTransactions[i].Comment, got.Comment)
				}

				if got.UserID != tt.userID {
					t.Errorf("want UserID %d; got %d", tt.userID, got.UserID)
				}
			}
		})
	}
}
