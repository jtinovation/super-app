package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type majorRepository struct {
	db *gorm.DB
}

func NewMajorRepository(db *gorm.DB) domain.MajorRepository {
	return &majorRepository{db: db}
}

func (r *majorRepository) FindAll(params dto.QueryParams) (*[]domain.Major, int64, error) {
	var majors []domain.Major
	var totalRows int64

	query := r.db.Model(&domain.Major{})

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where("LOWER(code) LIKE ? OR LOWER(name) LIKE ?", searchQuery, searchQuery)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&majors).Error; err != nil {
		return nil, 0, err
	}

	return &majors, totalRows, nil
}

func (r *majorRepository) FindByID(id string) (*domain.Major, error) {
	var major domain.Major
	if err := r.db.First(&major, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &major, nil
}

func (r *majorRepository) FindAllAsOptions() (*[]domain.Major, error) {
	var majors []domain.Major
	if err := r.db.Select("id", "name").Find(&majors).Error; err != nil {
		return nil, err
	}
	return &majors, nil
}

func (r *majorRepository) Create(major *domain.Major) (*domain.Major, error) {
	if err := r.db.Create(major).Error; err != nil {
		return nil, err
	}
	return major, nil
}

func (r *majorRepository) Update(id string, major *domain.Major) (*domain.Major, error) {
	var existingMajor domain.Major
	if err := r.db.First(&existingMajor, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingMajor).Updates(major).Error; err != nil {
		return nil, err
	}
	return &existingMajor, nil
}

func (r *majorRepository) Delete(id string) error {
	if err := r.db.First(&domain.Major{}, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Major{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
