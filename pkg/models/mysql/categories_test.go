package mysql

import "testing"

func TestCategoryInsert(t *testing.T) {
	t.Run("successful insert", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		categories := &CategoryModel{db}
		id, err := categories.Insert("Food", 1, 1, false)
		if err != nil {
			t.Fatal(err)
		}

		all, err := categories.All(1)
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

		if category.CategoryTypeID != 1 {
			t.Errorf("want CategoryTypeID %d; got %d", 1, category.CategoryTypeID)
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
	t.Run("successful fetch for specific user", func(t *testing.T) {
		db, teardown := newTestDB(t)
		defer teardown()

		categories := &CategoryModel{db}

		_, err := categories.Insert("Food", 1, 1, true)
		if err != nil {
			t.Fatal(err)
		}

		id, err := categories.Insert("Clothes", 1, 2, false)
		if err != nil {
			t.Fatal(err)
		}

		all, err := categories.All(2)
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

		if category.Name != "Clothes" {
			t.Errorf("want Name '%s'; got '%s'", "Clothes", category.Name)
		}

		if category.CategoryTypeID != 1 {
			t.Errorf("want CategoryTypeID %d; got %d", 1, category.CategoryTypeID)
		}

		if category.UserID != 2 {
			t.Errorf("want UserID %d; got %d", 2, category.UserID)
		}

		if category.Inactive != false {
			t.Errorf("want Inactive %t; got %t", false, category.Inactive)
		}
	})
}
