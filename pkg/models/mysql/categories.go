package mysql

import (
	"database/sql"

	"github.com/ck3g/homebugh-api/pkg/models"
)

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) Insert(name string, categoryType models.CategoryType, userID int64, inactive bool) (int64, error) {
	var id int64

	stmt := `INSERT INTO categories (name, category_type_id, user_id, inactive, updated_at)
	VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`

	res, err := m.DB.Exec(stmt, name, categoryType.ID, userID, inactive)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, nil
}

func (m *CategoryModel) All(userID int64, filters models.Filters) ([]*models.Category, models.Metadata, error) {
	categories := []*models.Category{}

	stmt := `SELECT c.id, c.name, c.category_type_id, ct.name, c.user_id, c.inactive, c.updated_at
	FROM categories AS c
	INNER JOIN category_types AS ct ON c.category_type_id = ct.id
	WHERE c.user_id = ?
	ORDER BY c.id
	LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(stmt, userID, filters.Limit(), filters.Offset())
	if err != nil {
		return categories, models.Metadata{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CategoryType.ID,
			&category.CategoryType.Name,
			&category.UserID,
			&category.Inactive,
			&category.UpdatedAt,
		)
		if err != nil {
			return categories, models.Metadata{}, err
		}

		categories = append(categories, &category)
	}

	totalRecords := 0
	countStmt := `SELECT COUNT(c.id)
	FROM categories AS c
	INNER JOIN category_types AS ct ON c.category_type_id = ct.id
	WHERE user_id = ?`
	err = m.DB.QueryRow(countStmt, userID).Scan(&totalRecords)
	if err != nil {
		return categories, models.Metadata{}, err
	}

	metadata := models.CalculateMetadata(totalRecords, filters.CurrentPage(), filters.Limit())

	return categories, metadata, nil
}
