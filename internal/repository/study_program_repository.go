package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type studyProgramRepository struct {
	db *gorm.DB
}

func NewStudyProgramRepository(db *gorm.DB) domain.StudyProgramRepository {
	return &studyProgramRepository{db: db}
}

func (r *studyProgramRepository) FindAll(params dto.QueryParams, majorId string) (*[]domain.StudyProgram, int64, error) {
	var studyPrograms []domain.StudyProgram
	var totalRows int64

	query := r.db.Model(&domain.StudyProgram{}).
		Select("m_study_program.*, m_major.name as major_name").
		Joins("LEFT JOIN m_major ON m_major.id = m_study_program.m_major_id")

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where(
			r.db.
				Where("LOWER(m_study_program.code) LIKE ?", searchQuery).
				Or("LOWER(m_study_program.name) LIKE ?", searchQuery).
				Or("LOWER(m_major.name) LIKE ?", searchQuery),
		)
	}

	if majorId != "" {
		query = query.Where("m_study_program.m_major_id = ?", majorId)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("major_name asc").Order("m_study_program.name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&studyPrograms).Error; err != nil {
		return nil, 0, err
	}

	return &studyPrograms, totalRows, nil
}

func (r *studyProgramRepository) FindByID(id string) (*domain.StudyProgram, error) {
	var studyProgram domain.StudyProgram
	if err := r.db.First(&studyProgram, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &studyProgram, nil
}

func (r *studyProgramRepository) FindAllAsOptions() (*[]domain.StudyProgram, error) {
	var studyPrograms []domain.StudyProgram
	if err := r.db.Select("id", "name").Find(&studyPrograms).Error; err != nil {
		return nil, err
	}
	return &studyPrograms, nil
}

func (r *studyProgramRepository) Create(studyProgram *domain.StudyProgram) (*domain.StudyProgram, error) {
	if err := r.db.Create(studyProgram).Error; err != nil {
		return nil, err
	}
	return studyProgram, nil
}

func (r *studyProgramRepository) Update(id string, studyProgram *domain.StudyProgram) (*domain.StudyProgram, error) {
	var existingStudyProgram domain.StudyProgram
	if err := r.db.First(&existingStudyProgram, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingStudyProgram).Updates(studyProgram).Error; err != nil {
		return nil, err
	}
	return &existingStudyProgram, nil
}

func (r *studyProgramRepository) Delete(id string) error {
	if err := r.db.First(&domain.StudyProgram{}, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&domain.StudyProgram{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
