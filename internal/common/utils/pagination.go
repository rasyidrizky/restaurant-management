package utils

import (
	"math"

	"project-2026-06-misoastory-be-go/internal/common/dto"

	"gorm.io/gorm"
)

// Paginate creates a GORM scope for pagination
func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 10
		}
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// CalculateMeta generates the pagination metadata
func CalculateMeta(total int64, page, limit int) dto.Meta {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return dto.Meta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
