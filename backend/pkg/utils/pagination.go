package utils

import (
	"math"
	"strconv"

	"github.com/saint0x/file-storage-app/backend/pkg/errors"
)

// PaginationParams represents the parameters for pagination
type PaginationParams struct {
	Page     int
	PageSize int
}

// Pagination represents the pagination information
type Pagination struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasPrevious bool `json:"has_previous"`
	HasNext     bool `json:"has_next"`
}

// NewPaginationFromRequest creates a new PaginationParams from request query parameters
func NewPaginationFromRequest(pageStr, pageSizeStr string) (*PaginationParams, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return nil, errors.BadRequest("Invalid page number")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		return nil, errors.BadRequest("Invalid page size")
	}

	return &PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CalculatePagination calculates pagination information
func CalculatePagination(totalItems, page, pageSize int) *Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return &Pagination{
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasPrevious: page > 1,
		HasNext:     page < totalPages,
	}
}

// CalculateOffset calculates the offset for SQL queries
func (p *PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}
