package mock

import "github.com/ck3g/homebugh-api/pkg/models"

type AccountModel struct{}

var (
	euro = models.Currency{ID: 1, Name: "Euro", Unit: "€"}

	user1BankAccount = &models.Account{
		ID:            1,
		Name:          "Bank",
		Balance:       1000,
		Currency:      euro,
		Status:        "active",
		ShowInSummary: true,
	}
	user1CashAccount = &models.Account{
		ID:            2,
		Name:          "Cash",
		Balance:       100.5,
		Currency:      euro,
		Status:        "active",
		ShowInSummary: true,
	}
	user2BankAccount = &models.Account{
		ID:            3,
		Name:          "Bank",
		Balance:       500,
		Currency:      euro,
		Status:        "active",
		ShowInSummary: true,
	}
	user2CashAccount = &models.Account{
		ID:            4,
		Name:          "Cash",
		Balance:       30.5,
		Currency:      euro,
		Status:        "active",
		ShowInSummary: true,
	}
)

func (m *AccountModel) Insert(name string, userID int64, currencyID int64, status string, showInSummary bool) (int64, error) {
	return 5, nil
}

func (m *AccountModel) All(userID int64, filters models.Filters) ([]*models.Account, models.Metadata, error) {
	var accounts []*models.Account

	switch userID {
	case 1:
		accounts = []*models.Account{user1BankAccount, user1CashAccount}
	case 2:
		accounts = []*models.Account{user2BankAccount, user2CashAccount}
	default:
		accounts = []*models.Account{}
	}

	metadata := models.CalculateMetadata(2, filters.CurrentPage(), filters.Limit())

	return accounts, metadata, nil
}
