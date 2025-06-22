package utils

import (
	"jti-super-app-go/internal/dto"

	"gorm.io/gorm"
)

func PaginateScope(pagination *dto.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}