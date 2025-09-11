package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Lab struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	MajorID   string `gorm:"column:m_major_id;type:char(36);not null"`
	Code      string `gorm:"type:varchar(255);not null"`
	Name      string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Major       Major         `gorm:"foreignKey:MajorID;references:ID"`
	EmployeeLab []EmployeeLab `gorm:"foreignKey:LabID;references:ID"`

	MajorName string `json:"major_name" gorm:"column:major_name;<-:false;->"`
}

func (Lab) TableName() string {
	return "m_lab"
}

type LabRepository interface {
	FindAll(params dto.QueryParams, majorId string) (*[]Lab, int64, error)
	FindByID(id string) (*Lab, error)
	FindAllAsOptions(majorId string) (*[]Lab, error)
	Create(lab *Lab) (*Lab, error)
	Update(id string, lab *Lab) (*Lab, error)
	Delete(id string) error
}
