package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Major struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	Code      string `gorm:"type:varchar(255);not null"`
	Name      string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	StudyPrograms []StudyProgram `gorm:"foreignKey:MajorID;references:ID"`
	// Employees     []Employee     `gorm:"foreignKey:MajorID"`
}

func (Major) TableName() string {
	return "m_major"
}

type MajorRepository interface {
	FindAll(params dto.QueryParams) (*[]Major, int64, error)
	FindByID(id string) (*Major, error)
	FindAllAsOptions() (*[]Major, error)
	Create(major *Major) (*Major, error)
	Update(id string, major *Major) (*Major, error)
	Delete(id string) error
}
