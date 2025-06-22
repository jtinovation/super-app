package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Class struct {
	ID             string `gorm:"type:char(36);primaryKey"`
	StudyProgramID string `gorm:"column:m_study_program_id;type:char(36);not null"`
	Code           string `gorm:"type:varchar(255);not null"`
	Name           string `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	StudyProgram StudyProgram `gorm:"foreignKey:StudyProgramID;references:ID"`

	StudyProgramName string `json:"study_program_name" gorm:"column:study_program_name;<-:false;->"`
	MajorName        string `json:"major_name" gorm:"column:major_name;<-:false;->"`
	MajorID          string `json:"major_id" gorm:"column:m_major_id;<-:false;->"`
}

func (Class) TableName() string {
	return "m_class"
}

type ClassRepository interface {
	FindAll(params dto.QueryParams, studyProgramId string, majorId string) (*[]Class, int64, error)
	FindByID(id string) (*Class, error)
	FindAllAsOptions() (*[]Class, error)
	Create(class *Class) (*Class, error)
	Update(id string, class *Class) (*Class, error)
	Delete(id string) error
}
