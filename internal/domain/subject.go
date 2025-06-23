package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Subject struct {
	ID             string `gorm:"type:char(36);primaryKey"`
	StudyProgramID string `gorm:"column:m_study_program_id;type:char(36);not null"`
	Code           string `gorm:"type:varchar(255);not null"`
	Name           string `gorm:"type:varchar(255);not null"`
	Status         string `gorm:"type:enum('ACTIVE','INACTIVE');default:'ACTIVE'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	StudyProgram     StudyProgram `gorm:"foreignKey:StudyProgramID"`
	Semesters        []Semester   `gorm:"many2many:m_subject_semester;foreignKey:ID;joinForeignKey:m_subject_id;References:ID;joinReferences:m_semester_id"`
	StudyProgramName string       `gorm:"column:study_program_name;<-:false;->"`
	StudyProgramID2  string       `gorm:"column:study_program_id;<-:false;->"`
}

func (Subject) TableName() string {
	return "m_subject"
}

type SubjectRepository interface {
	FindAll(params dto.QueryParams) (*[]Subject, int64, error)
	FindAllAsOptions(studyProgramID, semesterID string) (*[]Subject, error)
	FindByID(id string) (*Subject, error)
	Create(subject *Subject) (*Subject, error)
	Update(id string, subject *Subject) (*Subject, error)
	Delete(id string) error
}
