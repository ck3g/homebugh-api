package mock

import (
	"time"

	"github.com/ck3g/homebugh-api/pkg/models"
)

type TransactionModel struct{}

var (
	createdAt1, _           = time.Parse(time.RFC3339, "2022-10-31T16:24:05Z")
	createdAt2, _           = time.Parse(time.RFC3339, "2022-10-31T16:25:10Z")
	user1ExpenseTransaction = &models.Transaction{
		ID:        1,
		Amount:    20.0,
		Comment:   "food",
		UserID:    1,
		Category:  *foodCategory,
		Account:   *user1CashAccount,
		CreatedAt: &createdAt1,
	}
	user2ExpenseTransaction = &models.Transaction{
		ID:        2,
		Amount:    5.5,
		Comment:   "food 2",
		UserID:    2,
		Category:  *secondUserFoodCategory,
		Account:   *user2CashAccount,
		CreatedAt: &createdAt2,
	}
)

func (m *TransactionModel) Insert(amount float64, comment string, userID, categoryID, accountID int64) (int64, error) {
	return 3, nil
}

func (m *TransactionModel) All(userID int64, filters models.Filters) ([]*models.Transaction, models.Metadata, error) {
	var transactions []*models.Transaction

	switch userID {
	case 1:
		transactions = []*models.Transaction{user1ExpenseTransaction}
	case 2:
		transactions = []*models.Transaction{user2ExpenseTransaction}
	default:
		transactions = []*models.Transaction{}
	}

	metadata := models.CalculateMetadata(1, filters.CurrentPage(), filters.Limit())

	return transactions, metadata, nil
}
