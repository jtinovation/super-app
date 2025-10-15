package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectLecture struct {
	ID                string    `gorm:"type:char(36);primaryKey" json:"id"`
	SubjectSemesterID string    `gorm:"column:m_subject_semester_id;type:char(36);not null" json:"subject_semester_id"`
	EmployeeID        string    `gorm:"column:m_employee_id;type:char(36);not null" json:"employee_id"`
	CreatedAt         time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime:milli" json:"updated_at"`

	UserID                 string          `gorm:"column:user_id;<-:false;->" json:"user_id"`
	MajorIDEmployee        string          `gorm:"column:major_id_employee;<-:false;->" json:"major_id_employee"`
	StudyProgramIDEmployee string          `gorm:"column:study_program_id_employee;<-:false;->" json:"study_program_id_employee"`
	Employee               Employee        `gorm:"foreignKey:EmployeeID" json:"employee"`
	SubjectSemester        SubjectSemester `gorm:"foreignKey:SubjectSemesterID" json:"subject_semester"`
}

func (SubjectLecture) TableName() string {
	return "m_subject_lecture"
}

func (sl *SubjectLecture) BeforeCreate(tx *gorm.DB) (err error) {
	if sl.ID == "" {
		sl.ID = uuid.NewString()
	}
	return
}

type SubjectLectureRepository interface {
	FindAll(params dto.QueryParams) (*[]SubjectLecture, int64, error)
}
