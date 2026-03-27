package utils

import "math"

type PageMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NormalizePagination(page, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	return page, limit, offset
}

func BuildMeta(page, limit int, total int64) PageMeta {
	tp := int(math.Ceil(float64(total) / float64(limit)))
	if tp < 1 {
		tp = 1
	}
	return PageMeta{Page: page, Limit: limit, Total: total, TotalPages: tp}
}
