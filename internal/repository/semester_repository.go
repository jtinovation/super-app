package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type semesterRepository struct {
	db *gorm.DB
}

func NewSemesterRepository(db *gorm.DB) domain.SemesterRepository {
	return &semesterRepository{db: db}
}

func (r *semesterRepository) FindAll(params dto.QueryParams) (*[]domain.Semester, int64, error) {
	var semesters []domain.Semester
	var totalRows int64

	query := r.db.Model(&domain.Semester{}).
		Select("m_semester.*, m_session.session as session_name").
		Joins("LEFT JOIN m_session ON m_session.id = m_semester.m_session_id")

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where("LOWER(m_semester.year) LIKE ? OR LOWER(m_semester.semester) LIKE ?", searchQuery, searchQuery)
	}

	if sessionID, ok := params.Filter["session_id"]; ok && sessionID != "" {
		query = query.Where("m_semester.m_session_id = ?", sessionID)
	}

	if year, ok := params.Filter["year"]; ok && year != "" {
		query = query.Where("m_semester.year = ?", year)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		var sortOrder string
		if params.Sort == "session.session" {
			sortOrder = fmt.Sprintf("m_session.session %s", params.Order)
		} else {
			sortOrder = fmt.Sprintf("%s %s", params.Sort, params.Order)
		}
		query = query.Order(sortOrder)
	} else {
		query = query.Order("m_semester.year desc, m_semester.semester desc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&semesters).Error; err != nil {
		return nil, 0, err
	}
	return &semesters, totalRows, nil
}

func (r *semesterRepository) FindByID(id string) (*domain.Semester, error) {
	var semester domain.Semester
	if err := r.db.Preload("Session").First(&semester, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &semester, nil
}

func (r *semesterRepository) FindAllAsOptions(sessionID string) (*[]domain.Semester, error) {
	var semesters []domain.Semester
	query := r.db.Select("id", "year", "semester")
	if sessionID != "" {
		query = query.Where("m_session_id = ?", sessionID)
	}
	if err := query.Find(&semesters).Error; err != nil {
		return nil, err
	}
	return &semesters, nil
}

func (r *semesterRepository) Create(semester *domain.Semester) (*domain.Semester, error) {
	if err := r.db.Create(semester).Error; err != nil {
		return nil, err
	}
	return semester, nil
}

func (r *semesterRepository) Update(id string, semester *domain.Semester) (*domain.Semester, error) {
	var existingSemester domain.Semester
	if err := r.db.First(&existingSemester, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&existingSemester).Updates(semester).Error; err != nil {
		return nil, err
	}
	return &existingSemester, nil
}

func (r *semesterRepository) Delete(id string) error {
	if err := r.db.First(&domain.Semester{}, "id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.Delete(&domain.Semester{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *semesterRepository) SettingSubjectSemester(semesterID string, subjectIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var semester domain.Semester
		if err := tx.First(&semester, "id = ?", semesterID).Error; err != nil {
			return err
		}

		if err := tx.Where("m_semester_id = ?", semesterID).Delete(&domain.SubjectSemester{}).Error; err != nil {
			return err
		}

		if len(subjectIDs) == 0 {
			return nil
		}

		pivots := []domain.SubjectSemester{}
		for _, subID := range subjectIDs {
			pivots = append(pivots, domain.SubjectSemester{
				SemesterID: semesterID,
				SubjectID:  subID,
			})
		}
		if err := tx.Create(&pivots).Error; err != nil {
			return err
		}

		return nil
	})
}
