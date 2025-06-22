package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) domain.StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) FindAll(params dto.QueryParams) (*[]domain.Student, int64, error) {
	var students []domain.Student
	var totalRows int64

	query := r.db.Model(&domain.Student{}).
		Select(
			"m_student.*",
			"m_user.name as name",
			"m_user.img_path",
			"m_user.img_name",
			"m_class.name as class_name",
			"m_study_program.name as study_program_name",
			"m_study_program.id as study_program_id",
			"m_major.name as major_name",
			"m_major.id as major_id",
		).
		Joins("LEFT JOIN m_user ON m_user.id = m_student.m_user_id").
		Joins("LEFT JOIN m_class ON m_class.id = m_student.m_class_id").
		Joins("LEFT JOIN m_study_program ON m_study_program.id = m_class.m_study_program_id").
		Joins("LEFT JOIN m_major ON m_major.id = m_study_program.m_major_id")

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where(
			"LOWER(m_user.name) LIKE ? OR LOWER(m_student.nim) LIKE ?",
			searchQuery, searchQuery,
		)
	}

	if majorID, ok := params.Filter["major_id"]; ok && majorID != "" {
		query = query.Where("m_major.id = ?", majorID)
	}
	if spID, ok := params.Filter["study_program_id"]; ok && spID != "" {
		query = query.Where("m_study_program.id = ?", spID)
	}
	if classID, ok := params.Filter["class_id"]; ok && classID != "" {
		query = query.Where("m_class.id = ?", classID)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("m_user.name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&students).Error; err != nil {
		return nil, 0, err
	}

	return &students, totalRows, nil
}

func (r *studentRepository) FindByID(id string) (*domain.Student, error) {
	var student domain.Student
	if err := r.db.Preload("User").Preload("Class").First(&student, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Create(student *domain.Student) (*domain.Student, error) {
	if err := r.db.Create(student).Error; err != nil {
		return nil, err
	}
	return student, nil
}

func (r *studentRepository) Update(id string, student *domain.Student) (*domain.Student, error) {
	panic("not implemented")
}

func (r *studentRepository) Delete(id string) error {
	panic("not implemented")
}
