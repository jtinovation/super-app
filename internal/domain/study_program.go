package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type StudyProgram struct {
	ID        string         `gorm:"type:char(36);primaryKey" json:"id"`
	MajorID   string         `gorm:"column:m_major_id;type:char(36);not null" json:"major_id"`
	Code      string         `gorm:"type:varchar(255);not null" json:"code"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Major Major `gorm:"foreignKey:MajorID;references:ID"`

	MajorName string `json:"major_name" gorm:"column:major_name;<-:false;->"`
}

func (StudyProgram) TableName() string {
	return "m_study_program"
}

type StudyProgramRepository interface {
	FindAll(params dto.QueryParams, majorId string) (*[]StudyProgram, int64, error)
	FindByID(id string) (*StudyProgram, error)
	FindAllAsOptions(majorId string) (*[]StudyProgram, error)
	Create(studyProgram *StudyProgram) (*StudyProgram, error)
	Update(id string, studyProgram *StudyProgram) (*StudyProgram, error)
	Delete(id string) error
}
