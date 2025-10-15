package repository

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"gorm.io/gorm"
)

type subjectLectureRepository struct {
	db *gorm.DB
}

func NewSubjectLectureRepository(db *gorm.DB) *subjectLectureRepository {
	return &subjectLectureRepository{db: db}
}

func (r *subjectLectureRepository) FindAll(params dto.QueryParams) (*[]domain.SubjectLecture, int64, error) {
	var subjectLectures []domain.SubjectLecture
	var totalRows int64

	query := r.db.Model(&domain.SubjectLecture{}).
		Select("m_subject_lecture.*", "m_employee.m_user_id as user_id", "m_employee.m_major_id as major_id_employee", "m_employee.m_study_program_id as study_program_id_employee").
		Joins("JOIN m_employee ON m_employee.id = m_subject_lecture.m_employee_id")

	if params.Search != "" {
		searchQuery := "%" + params.Search + "%"
		query = query.Joins("JOIN m_employee ON m_employee.id = m_subject_lecture.m_employee_id").
			Where("m_employee.name LIKE ?", searchQuery)
	}

	if params.Filter != nil {
		if employeeID, ok := params.Filter["employee_id"]; ok && employeeID != "" {
			query = query.Where("m_subject_lecture.m_employee_id = ?", employeeID)
		}

		if semesterID, ok := params.Filter["semester_id"]; ok && semesterID != "" {
			query = query.Preload("SubjectSemester", "m_semester_id = ?", semesterID)
		}

		if employeeIDs, ok := params.Filter["employee_ids"]; ok && employeeIDs != "" {
			query = query.Where("m_subject_lecture.m_employee_id IN ?", employeeIDs)
		}

		if userIDs, ok := params.Filter["user_ids"]; ok && userIDs != "" {
			query = query.Where("m_employee.m_user_id IN ?", userIDs)
		}

		if semesterIDs, ok := params.Filter["semester_ids"]; ok && semesterIDs != "" {
			query = query.Preload("SubjectSemester.Subject.StudyProgram").Joins("JOIN m_subject_semester ON m_subject_semester.id = m_subject_lecture.m_subject_semester_id").
				Where("m_subject_semester.m_semester_id IN ?", semesterIDs)
		}
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := params.Sort + " " + params.Order
		query = query.Order(sortOrder)
	} else {
		query = query.Order("m_subject_lecture.created_at desc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&subjectLectures).Error; err != nil {
		return nil, 0, err
	}

	return &subjectLectures, totalRows, nil
}
