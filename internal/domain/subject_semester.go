package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectSemester struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	SubjectID  string    `gorm:"column:m_subject_id;type:char(36);not null" json:"subject_id"`
	SemesterID string    `gorm:"column:m_semester_id;type:char(36);not null" json:"semester_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Subject         Subject          `gorm:"foreignKey:SubjectID" json:"subject"`
	Semester        Semester         `gorm:"foreignKey:SemesterID" json:"semester"`
	SubjectLectures []SubjectLecture `gorm:"foreignKey:SubjectSemesterID" json:"subject_lectures"`
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
