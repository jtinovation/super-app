package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Semester struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	SessionID string `gorm:"column:m_session_id;type:char(36);not null"`
	Year      int    `gorm:"type:int;not null"`
	Semester  string `gorm:"type:varchar(2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Session  Session   `gorm:"foreignKey:SessionID;references:ID"`
	Subjects []Subject `gorm:"many2many:m_subject_semester;foreignKey:ID;joinForeignKey:m_semester_id;References:ID;joinReferences:m_subject_id"`

	SessionName string `gorm:"column:session_name;<-:false;->"`
}

func (Semester) TableName() string {
	return "m_semester"
}

type SemesterRepository interface {
	FindAll(params dto.QueryParams) (*[]Semester, int64, error)
	FindByID(id string) (*Semester, error)
	FindAllAsOptions(sessionID string) (*[]Semester, error)
	Create(semester *Semester) (*Semester, error)
	Update(id string, semester *Semester) (*Semester, error)
	Delete(id string) error
	SettingSubjectSemester(semesterID string, subjectIDs []string) error
}
