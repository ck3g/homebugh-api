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

	query := `SELECT a.id, a.name, a.funds, a.currency_id, c.name, c.unit, a.status, a.show_in_summary
	FROM accounts AS a
	INNER JOIN currencies AS c ON a.currency_id = c.id
	WHERE a.user_id = ?
	ORDER BY a.id
	LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(query, userID, filters.Limit(), filters.Offset())
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
			&account.Currency.ID,
			&account.Currency.Name,
			&account.Currency.Unit,
			&account.Status,
			&account.ShowInSummary,
		)
		if err != nil {
			return accounts, models.Metadata{}, err
		}

		accounts = append(accounts, &account)
	}

	totalRecords := 0
	countQuery := `SELECT COUNT(a.id)
	FROM accounts AS a
	INNER JOIN currencies AS c ON a.currency_id = c.id
	WHERE a.user_id = ?`
	err = m.DB.QueryRow(countQuery, userID).Scan(&totalRecords)
	if err != nil {
		return accounts, models.Metadata{}, err
	}

	metadata := models.CalculateMetadata(totalRecords, filters.CurrentPage(), filters.Limit())

	return accounts, metadata, nil
}
