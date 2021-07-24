package persistence

import (
	"math"

	"gorm.io/gorm"
)

type PaginationInput struct {
	Page  uint
	Limit uint
}

func (pi *PaginationInput) GetOffset() uint {
	return (pi.Page - 1) * pi.Limit
}

type PaginationOutput struct {
	Page       uint
	Limit      uint
	TotalItems uint
	TotalPages uint
}

func Paginate(value interface{}, db *gorm.DB, input *PaginationInput, output *PaginationOutput) func(db *gorm.DB) *gorm.DB {
	var totalItems int64
	db.Model(value).Count(&totalItems)

	output.Limit = input.Limit
	output.Page = input.Page
	output.TotalItems = uint(totalItems)
	totalPages := int(math.Ceil(float64(totalItems) / float64(input.Limit)))
	output.TotalPages = uint(totalPages)

	return func(db *gorm.DB) *gorm.DB {
		return db.
			Offset(int(input.GetOffset())).
			Limit(int(input.Limit))
	}
}
