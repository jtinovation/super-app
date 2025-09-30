package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectSemester struct {
	ID         string `gorm:"type:char(36);primaryKey"`
	SubjectID  string `gorm:"column:m_subject_id;type:char(36);not null"`
	SemesterID string `gorm:"column:m_semester_id;type:char(36);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Subject         Subject          `gorm:"foreignKey:SubjectID"`
	Semester        Semester         `gorm:"foreignKey:SemesterID"`
	SubjectLectures []SubjectLecture `gorm:"foreignKey:SubjectSemesterID"`
	// Lecturers []Employee `gorm:"many2many:m_subject_lecture;using:jti-super-app-go/internal/domain.SubjectLecture;foreignKey:ID;joinForeignKey:m_subject_semester_id;References:ID;joinReferences:m_employee_id"`
}

func (SubjectSemester) TableName() string {
	return "m_subject_semester"
}

func (ss *SubjectSemester) BeforeCreate(tx *gorm.DB) (err error) {
	if ss.ID == "" {
		ss.ID = uuid.NewString()
	}
	return
}

type SubjectSemesterRepository interface {
	GetLectureOnSubject(studyProgramID, semesterID string) (*[]SubjectSemester, error)
	StoreLectureOnSubject(data []dto.LectureMappingDTO) error
}
