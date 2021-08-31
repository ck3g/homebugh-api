package mysql

import (
	"testing"

	"github.com/ck3g/homebugh-api/pkg/models"
)

func TestCategoryInsert(t *testing.T) {
	t.Run("successful insert", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		categories := &CategoryModel{db}
		expense := models.CategoryTypeID(2)
		id, err := categories.Insert("Food", expense, 1, false)
		if err != nil {
			t.Fatal(err)
		}

		filters := models.Filters{
			Page:     1,
			PageSize: 20,
		}

		all, _, err := categories.All(1, filters)
		if err != nil {
			t.Fatal(err)
		}

		if len(all) != 1 {
			t.Errorf("want 1 category; got %d", len(all))
		}

		category := all[0]
		if category.ID != id {
			t.Errorf("want ID %d; got %d", id, category.ID)
		}

		if category.Name != "Food" {
			t.Errorf("want Name '%s'; got '%s'", "Food", category.Name)
		}

		if category.CategoryTypeID != expense {
			t.Errorf("want CategoryTypeID %d; got %d", expense, category.CategoryTypeID)
		}

		if category.UserID != 1 {
			t.Errorf("want UserID %d; got %d", 1, category.UserID)
		}

		if category.Inactive != false {
			t.Errorf("want Inactive %t; got %t", false, category.Inactive)
		}
	})
}

func TestCategoryAll(t *testing.T) {
	db, teardown := newTestDB(t)
	defer teardown()

	categories := &CategoryModel{db}
	income := models.CategoryTypeID(1)
	expense := models.CategoryTypeID(2)

	allCategories := map[int]*models.Category{
		0: {Name: "Food", CategoryTypeID: expense, UserID: 1, Inactive: true},
		1: {Name: "Clothes", CategoryTypeID: expense, UserID: 2},
		2: {Name: "Salary", CategoryTypeID: income, UserID: 2},
	}
	ids := map[int]int64{}

	for i, c := range allCategories {
		id, err := categories.Insert(c.Name, c.CategoryTypeID, c.UserID, c.Inactive)
		if err != nil {
			t.Fatal(err)
		}

		ids[i] = id
	}

	tests := []struct {
		name             string
		wantCategories   []*models.Category
		wantCount        int
		wantTotalRecords int
		userID           int64
		filters          models.Filters
	}{
		{
			name: "successful fetch for specific user",
			wantCategories: []*models.Category{
				{ID: ids[1], Name: "Clothes"},
				{ID: ids[2], Name: "Salary"},
			},
			wantCount:        2,
			wantTotalRecords: 2,
			userID:           2,
			filters:          models.Filters{Page: 1, PageSize: 20},
		},
		{
			name: "categories from the first page",
			wantCategories: []*models.Category{
				{ID: ids[1], Name: "Clothes"},
			},
			wantCount:        1,
			wantTotalRecords: 2,
			userID:           2,
			filters:          models.Filters{Page: 1, PageSize: 1},
		},
		{
			name: "categories from the second page",
			wantCategories: []*models.Category{
				{ID: ids[2], Name: "Salary"},
			},
			wantCount:        1,
			wantTotalRecords: 2,
			userID:           2,
			filters:          models.Filters{Page: 2, PageSize: 1},
		},
		{
			name:             "categories from the greater than last page",
			wantCategories:   []*models.Category{},
			wantCount:        0,
			wantTotalRecords: 2,
			userID:           2,
			filters:          models.Filters{Page: 100, PageSize: 1},
		},
		{
			name: "categories from the negative page",
			wantCategories: []*models.Category{
				{ID: ids[1], Name: "Clothes"},
			},
			wantCount:        1,
			wantTotalRecords: 2,
			userID:           2,
			filters:          models.Filters{Page: -1, PageSize: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			all, metadata, err := categories.All(tt.userID, tt.filters)
			if err != nil {
				t.Fatal(err)
			}

			if len(all) != tt.wantCount {
				t.Errorf("want %d categories; got %d", tt.wantCount, len(all))
			}

			if metadata.TotalRecords != tt.wantTotalRecords {
				t.Errorf("want %d total records; got %d", tt.wantTotalRecords, metadata.TotalRecords)
			}

			for i, got := range all {
				if got.ID != tt.wantCategories[i].ID {
					t.Errorf("want ID %d; got %d", tt.wantCategories[i].ID, got.ID)
				}

				if got.Name != tt.wantCategories[i].Name {
					t.Errorf("want Name '%s'; got '%s'", tt.wantCategories[i].Name, got.Name)
				}

				if got.UserID != tt.userID {
					t.Errorf("want UserID %d; got %d", tt.userID, got.UserID)
				}
			}
		})
	}
}
