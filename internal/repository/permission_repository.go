package repository

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *permissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) FindByID(id string) (*domain.Permission, error) {
	var permission domain.Permission
	if err := r.db.First(&permission, "uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindAll(params dto.QueryParams) (*[]domain.Permission, int64, error) {
	var permissions []domain.Permission
	var totalRows int64

	query := r.db.Model(&domain.Permission{})

	if params.Search != "" {
		searchQuery := "%" + params.Search + "%"
		query = query.Where("name LIKE ? OR uuid LIKE ?", searchQuery, searchQuery)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := params.Sort + " " + params.Order
		query = query.Order(sortOrder)
	} else {
		query = query.Order("name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&permissions).Error; err != nil {
		return nil, 0, err
	}
	return &permissions, totalRows, nil
}

func (r *permissionRepository) Create(permission *domain.Permission) (*domain.Permission, error) {
	if err := r.db.Create(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

func (r *permissionRepository) Update(id string, permission *domain.Permission) (*domain.Permission, error) {
	if err := r.db.Model(&domain.Permission{}).Where("uuid = ?", id).Updates(permission).Error; err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *permissionRepository) Delete(id string) error {
	if err := r.db.Delete(&domain.Permission{}, "uuid = ?", id).Error; err != nil {
		return err
	}
	return nil
}
