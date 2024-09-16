package pagination

import (
	"net/http"
	"strconv"
)

type PaginationRequest struct {
	PageSize int    `json:"pageSize"`
	Page     int    `json:"page"`
	SortBy   string `json:"sortBy"`
	Order    string `json:"order"`
}

func (r PaginationRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

func NewPaginationRequest(r *http.Request) *PaginationRequest {
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize == 0 || err != nil {
		pageSize = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 || err != nil {
		page = 1
	}

	sortBy := r.URL.Query().Get("page")
	if sortBy == "" {
		sortBy = "updated_at"
	}

	order := r.URL.Query().Get("order")
	if order == "" {
		order = "desc"
	}

	return &PaginationRequest{
		PageSize: pageSize,
		Page:     page,
		SortBy:   sortBy,
		Order:    order,
	}
}
