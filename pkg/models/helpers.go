package models

import "math"

const defaultPageSize = 20

type Filters struct {
	Page     int
	PageSize int
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func (f Filters) Limit() int {
	if f.PageSize < 1 || f.PageSize > 50 {
		return defaultPageSize
	}

	return f.PageSize
}

func (f Filters) Offset() int {
	if f.Page < 1 {
		f.Page = 1
	}

	return (f.Page - 1) * f.Limit()
}

func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}
