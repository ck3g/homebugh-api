package mysql

import (
	"database/sql"

	"github.com/ck3g/homebugh-api/pkg/models"
)

type TransactionModel struct {
	DB *sql.DB
}

func (m *TransactionModel) Insert(amount float64, comment string, userID, categoryID, accountID int64) (int64, error) {
	var id int64

	query := `INSERT INTO transactions (summ, comment, user_id, category_id, account_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	args := []interface{}{amount, comment, userID, categoryID, accountID}
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

func (m *TransactionModel) All(userID int64, filters models.Filters) ([]*models.Transaction, models.Metadata, error) {
	transactions := []*models.Transaction{}

	query := `SELECT
		t.id,
		t.summ,
		t.comment,
		t.user_id,
		t.category_id,
		c.name,
		c.inactive,
		c.category_type_id,
		ct.name,
		t.account_id,
		a.name,
		a.status,
		a.show_in_summary,
		a.currency_id,
		cr.name,
		cr.unit
	FROM transactions AS t
	INNER JOIN categories AS c ON t.category_id = c.id
	INNER JOIN category_types AS ct ON c.category_type_id = ct.id
	INNER JOIN accounts AS a ON t.account_id = a.id
	INNER JOIN currencies AS cr ON a.currency_id = cr.id
	WHERE t.user_id = ?
	ORDER BY t.id DESC
	LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(query, userID, filters.Limit(), filters.Offset())
	if err != nil {
		return transactions, models.Metadata{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.Comment,
			&transaction.UserID,
			&transaction.Category.ID,
			&transaction.Category.Name,
			&transaction.Category.Inactive,
			&transaction.Category.CategoryType.ID,
			&transaction.Category.CategoryType.Name,
			&transaction.Account.ID,
			&transaction.Account.Name,
			&transaction.Account.Status,
			&transaction.Account.ShowInSummary,
			&transaction.Account.Currency.ID,
			&transaction.Account.Currency.Name,
			&transaction.Account.Currency.Unit,
		)
		if err != nil {
			return transactions, models.Metadata{}, err
		}

		transactions = append(transactions, &transaction)
	}

	totalRecords := 0
	countQuery := `SELECT COUNT(t.id)
	FROM transactions AS t
	WHERE t.user_id = ?`
	err = m.DB.QueryRow(countQuery, userID).Scan(&totalRecords)
	if err != nil {
		return transactions, models.Metadata{}, err
	}

	metadata := models.CalculateMetadata(totalRecords, filters.CurrentPage(), filters.Limit())

	return transactions, metadata, nil
}
