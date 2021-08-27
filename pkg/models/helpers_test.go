package models

import "testing"

func TestFilterLimit(t *testing.T) {
	tests := []struct {
		name      string
		pageSize  int
		wantLimit int
	}{
		{"uses specified page size 1", 1, 1},
		{"uses specified page size 50", 50, 50},
		{"uses default page size when provided size less than 1", -1, 20},
		{"uses default page size when provided size more than 50", 51, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := Filters{PageSize: tt.pageSize}
			limit := filters.Limit()

			if limit != tt.wantLimit {
				t.Errorf("want limit %d; got %d", tt.wantLimit, limit)
			}
		})
	}
}

func TestFilterOffset(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantOffset int
	}{
		{"with page 1", 1, 20, 0},
		{"with page 0", 0, 20, 0},
		{"with page 2", 2, 20, 20},
		{"with page 3 and size 2", 3, 2, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := Filters{Page: tt.page, PageSize: tt.pageSize}
			offset := filters.Offset()

			if offset != tt.wantOffset {
				t.Errorf("want offset %d; got %d", tt.wantOffset, offset)
			}
		})
	}
}
