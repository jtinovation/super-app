package repository

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByID(id string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.First(&role, "uuid = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll(params dto.QueryParams) (*[]domain.Role, int64, error) {
	var roles []domain.Role
	var totalRows int64

	query := r.db.Model(&domain.Role{}).Preload("Permissions")

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

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return &roles, totalRows, nil
}

func (r *roleRepository) Create(role *domain.Role) (*domain.Role, error) {
	if err := r.db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) Update(id string, role *domain.Role) (*domain.Role, error) {
	// Start a transaction to ensure atomicity
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// Update the role fields
	if err := tx.Model(&domain.Role{}).Where("uuid = ?", id).Updates(role).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update the permissions association if Permissions is not nil
	if role.Permissions != nil {
		var updatedRole domain.Role
		if err := tx.Where("uuid = ?", id).First(&updatedRole).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(&updatedRole).Association("Permissions").Replace(role.Permissions); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (r *roleRepository) Delete(id string) error {
	if err := r.db.Where("uuid = ?", id).Delete(&domain.Role{}).Error; err != nil {
		return err
	}
	return nil
}
