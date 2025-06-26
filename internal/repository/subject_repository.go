package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/pkg/constants"

	"gorm.io/gorm"
)

type subjectRepository struct{ db *gorm.DB }

func NewSubjectRepository(db *gorm.DB) domain.SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) FindAll(params dto.QueryParams) (*[]domain.Subject, int64, error) {
	var subjects []domain.Subject
	var totalRows int64

	query := r.db.Model(&domain.Subject{}).
		Select([]string{
			"m_subject.id",
			"m_subject.m_study_program_id as study_program_id",
			"m_subject.code",
			"m_subject.name",
			"m_subject.status",
			"m_study_program.name as study_program_name",
		}).
		Joins("LEFT JOIN m_study_program ON m_subject.m_study_program_id = m_study_program.id")

	if params.Search != "" {
		search := fmt.Sprintf("%%%s%%", params.Search)
		query = query.Where(
			r.db.Where("m_subject.code LIKE ?", search).
				Or("m_subject.name LIKE ?", search).
				Or("m_study_program.name LIKE ?", search),
		)
	}

	if spID, ok := params.Filter["study_program_id"]; ok && spID != "" {
		query = query.Where("m_subject.m_study_program_id = ?", spID)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		query = query.Order(fmt.Sprintf("%s %s", params.Sort, params.Order))
	} else {
		query = query.Order("m_subject.name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&subjects).Error; err != nil {
		return nil, 0, err
	}
	return &subjects, totalRows, nil
}

func (r *subjectRepository) FindAllAsOptions(studyProgramID, semesterID string) (*[]domain.Subject, error) {
	var subjects []domain.Subject

	query := r.db.Table("m_subject").
		Select("MIN(m_subject.id) as id, m_subject.code, MIN(m_subject.name) as name").
		Where("m_subject.status = ?", constants.StatusActive).
		Where("m_subject.deleted_at IS NULL")

	if studyProgramID != "" {
		query = query.Where("m_subject.m_study_program_id = ?", studyProgramID)
	}

	if semesterID != "" {
		query = query.Joins("JOIN m_subject_semester ON m_subject_semester.m_subject_id = m_subject.id").
			Where("m_subject_semester.m_semester_id = ?", semesterID)
	}

	err := query.
		Group("m_subject.code").
		Order("name asc").
		Scan(&subjects).Error

	return &subjects, err
}

func (r *subjectRepository) FindByID(id string) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.db.First(&subject, "id = ?", id).Error
	return &subject, err
}

func (r *subjectRepository) Create(subject *domain.Subject) (*domain.Subject, error) {
	if err := r.db.Create(subject).Error; err != nil {
		return nil, err
	}
	return subject, nil
}

func (r *subjectRepository) Update(id string, subject *domain.Subject) (*domain.Subject, error) {
	if err := r.db.First(&domain.Subject{}, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&domain.Subject{ID: id}).Updates(subject).Error; err != nil {
		return nil, err
	}
	return subject, nil
}

func (r *subjectRepository) Delete(id string) error {
	if err := r.db.First(&domain.Subject{}, "id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Delete(&domain.Subject{}, "id = ?", id).Error
}
