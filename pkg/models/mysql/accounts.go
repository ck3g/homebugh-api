package mysql

import (
	"database/sql"

	"github.com/ck3g/homebugh-api/pkg/models"
)

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) Insert(name string, userID int64, currencyID int64, status string, showInSummary bool) (int64, error) {
	var id int64

	query := `INSERT INTO accounts (name, user_id, currency_id, status, show_in_summary, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	args := []interface{}{name, userID, currencyID, status, showInSummary}
	res, err := m.DB.Exec(query, args...)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func (m *AccountModel) All(userID int64, filters models.Filters) ([]*models.Account, models.Metadata, error) {
	accounts := []*models.Account{}

	query := `SELECT id, name, funds, currency_id, status, show_in_summary
	FROM accounts
	WHERE user_id = ?
	ORDER BY id`

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return accounts, models.Metadata{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var account models.Account

		err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.Balance,
			&account.CurrencyID,
			&account.Status,
			&account.ShowInSummary,
		)
		if err != nil {
			return accounts, models.Metadata{}, err
		}

		accounts = append(accounts, &account)
	}

	metadata := models.CalculateMetadata(1, filters.CurrentPage(), filters.Limit())

	return accounts, metadata, nil
}
