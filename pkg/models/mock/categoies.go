package mock

import (
	"github.com/ck3g/homebugh-api/pkg/models"
)

var (
	foodCategory = &models.Category{
		ID:           1,
		Name:         "Food",
		CategoryType: models.CategoryType{ID: 2, Name: "expense"},
		UserID:       1,
		Inactive:     false,
		UpdatedAt:    &oneDayAgo,
	}
	secondUserFoodCategory = &models.Category{
		ID:           2,
		Name:         "Groceries",
		CategoryType: models.CategoryType{ID: 2, Name: "expense"},
		UserID:       2,
		Inactive:     false,
		UpdatedAt:    &oneDayAgo,
	}
)

type CategoryModel struct{}

func (m *CategoryModel) Insert(name string, categoryType models.CategoryType, userID int64, inactive bool) (int64, error) {
	return 3, nil
}

func (m *CategoryModel) All(userID int64, filters models.Filters) ([]*models.Category, models.Metadata, error) {
	var categories []*models.Category

	switch userID {
	case 1:
		categories = []*models.Category{foodCategory}
	case 2:
		categories = []*models.Category{secondUserFoodCategory}
	default:
		categories = []*models.Category{}
	}

	metadata := models.CalculateMetadata(1, filters.CurrentPage(), filters.Limit())

	return categories, metadata, nil
}
