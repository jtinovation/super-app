package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) domain.ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) FindAll(params dto.QueryParams, studyProgramId string, majorId string) (*[]domain.Class, int64, error) {
	var classes []domain.Class
	var totalRows int64

	query := r.db.Model(&domain.Class{}).
		Select("m_class.*, m_study_program.name as study_program_name, m_major.name as major_name, m_major.id as m_major_id").
		Joins("LEFT JOIN m_study_program ON m_study_program.id = m_class.m_study_program_id").
		Joins("LEFT JOIN m_major ON m_major.id = m_study_program.m_major_id")

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where(
			r.db.
				Where("LOWER(m_class.code) LIKE ?", searchQuery).
				Or("LOWER(m_class.name) LIKE ?", searchQuery).
				Or("LOWER(m_study_program.name) LIKE ?", searchQuery).
				Or("LOWER(m_major.name) LIKE ?", searchQuery),
		)
	}

	if studyProgramId != "" {
		query = query.Where("m_class.m_study_program_id = ?", studyProgramId)
	}

	if majorId != "" {
		query = query.Where("m_study_program.m_major_id = ?", majorId)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		var sortOrder string
		if params.Sort == "study_program.name" {
			sortOrder = fmt.Sprintf("m_study_program.name %s", params.Order)
		} else if params.Sort == "major.name" {
			sortOrder = fmt.Sprintf("m_major.name %s", params.Order)
		} else {
			sortOrder = fmt.Sprintf("%s %s", params.Sort, params.Order)
		}
		query = query.Order(sortOrder)
	} else {
		query = query.Order("m_major.name asc").
			Order("m_study_program.name asc").
			Order("m_class.name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return &classes, totalRows, nil
}

func (r *classRepository) FindByID(id string) (*domain.Class, error) {
	var class domain.Class
	if err := r.db.First(&class, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) FindAllAsOptions(studyProgramId string) (*[]domain.Class, error) {
	var classs []domain.Class
	query := r.db.Select("id", "name")
	if studyProgramId != "" {
		query = query.Where("m_study_program_id = ?", studyProgramId)
	}

	if err := query.Find(&classs).Error; err != nil {
		return nil, err
	}
	return &classs, nil
}

func (r *classRepository) Create(class *domain.Class) (*domain.Class, error) {
	if err := r.db.Create(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}

func (r *classRepository) Update(id string, class *domain.Class) (*domain.Class, error) {
	var existingClass domain.Class
	if err := r.db.First(&existingClass, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingClass).Updates(class).Error; err != nil {
		return nil, err
	}
	return &existingClass, nil
}

func (r *classRepository) Delete(id string) error {
	if err := r.db.First(&domain.Class{}, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Class{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
