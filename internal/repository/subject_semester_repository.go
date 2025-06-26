package repository

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"gorm.io/gorm"
)

type subjectSemesterRepository struct{ db *gorm.DB }

func NewSubjectSemesterRepository(db *gorm.DB) domain.SubjectSemesterRepository {
	return &subjectSemesterRepository{db: db}
}

func (r *subjectSemesterRepository) GetLectureOnSubject(studyProgramID, semesterID string) (*[]domain.SubjectSemester, error) {
	var subjectSemesters []domain.SubjectSemester
	err := r.db.
		Preload("Subject").
		Preload("Lecturers.User").
		Joins("JOIN m_subject ON m_subject.id = m_subject_semester.m_subject_id").
		Where("m_subject.m_study_program_id = ?", studyProgramID).
		Where("m_subject_semester.m_semester_id = ?", semesterID).
		Order("m_subject.name ASC").
		Find(&subjectSemesters).Error
	return &subjectSemesters, err
}

func (r *subjectSemesterRepository) StoreLectureOnSubject(data []dto.LectureMappingDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range data {
			var subjectSemester domain.SubjectSemester
			if err := tx.First(&subjectSemester, "id = ?", item.SubjectSemesterID).Error; err != nil {
				return err
			}

			if err := tx.Where("m_subject_semester_id = ?", item.SubjectSemesterID).Delete(&domain.SubjectLecture{}).Error; err != nil {
				return err
			}

			if len(item.LectureIDs) == 0 {
				continue
			}

			pivots := []domain.SubjectLecture{}
			for _, lectureID := range item.LectureIDs {
				pivots = append(pivots, domain.SubjectLecture{
					SubjectSemesterID: item.SubjectSemesterID,
					EmployeeID:        lectureID,
				})
			}

			if err := tx.Create(&pivots).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
