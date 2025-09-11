package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type labRepository struct {
	db *gorm.DB
}

func NewLabRepository(db *gorm.DB) domain.LabRepository {
	return &labRepository{db: db}
}

func (r *labRepository) FindAll(params dto.QueryParams, majorId string) (*[]domain.Lab, int64, error) {
	var labs []domain.Lab
	var totalRows int64

	query := r.db.Model(&domain.Lab{}).
		Preload("EmployeeLab", "status = ?", "ACTIVE").
		Select("m_lab.*, m_major.name as major_name").
		Joins("LEFT JOIN m_major ON m_major.id = m_lab.m_major_id")

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where(
			r.db.
				Where("LOWER(m_lab.code) LIKE ?", searchQuery).
				Or("LOWER(m_lab.name) LIKE ?", searchQuery).
				Or("LOWER(m_major.name) LIKE ?", searchQuery),
		)
	}

	if majorId != "" {
		query = query.Where("m_lab.m_major_id = ?", majorId)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("major_name asc").Order("m_lab.name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&labs).Error; err != nil {
		return nil, 0, err
	}

	return &labs, totalRows, nil
}

func (r *labRepository) FindByID(id string) (*domain.Lab, error) {
	var lab domain.Lab
	if err := r.db.First(&lab, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lab, nil
}

func (r *labRepository) FindAllAsOptions(majorId string) (*[]domain.Lab, error) {
	var labs []domain.Lab
	query := r.db.Select("id", "name")
	if majorId != "" {
		query = query.Where("m_major_id = ?", majorId)
	}
	if err := query.Find(&labs).Error; err != nil {
		return nil, err
	}
	return &labs, nil
}

func (r *labRepository) Create(lab *domain.Lab) (*domain.Lab, error) {
	if err := r.db.Create(lab).Error; err != nil {
		return nil, err
	}
	return lab, nil
}

func (r *labRepository) Update(id string, lab *domain.Lab) (*domain.Lab, error) {
	var existingLab domain.Lab
	if err := r.db.First(&existingLab, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingLab).Updates(lab).Error; err != nil {
		return nil, err
	}
	return &existingLab, nil
}

func (r *labRepository) Delete(id string) error {
	if err := r.db.First(&domain.Lab{}, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Lab{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
