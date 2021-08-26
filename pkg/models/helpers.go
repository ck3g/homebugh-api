package models

const defaultPageSize = 20

type Filters struct {
	Page     int
	PageSize int
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
