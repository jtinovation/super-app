package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID        string  `gorm:"type:char(36);primaryKey"`
	UserID    string  `gorm:"column:m_user_id;type:char(36);not null"`
	MajorID   *string `gorm:"column:m_major_id;type:char(36)"` // Nullable
	Nip       string  `gorm:"type:varchar(255);not null"`
	Position  string  `gorm:"type:enum('LECTURER','STAFF');not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Major Major `gorm:"foreignKey:MajorID;references:ID"`

	// Fields from join for FindAll operation
	// Based on select('m_employee.*', 'm_user.name as name', ...) in EmployeeRepository.php
	Name    string `json:"name" gorm:"column:name;<-:false;->"`
	Email   string `json:"email" gorm:"column:email;<-:false;->"`
	ImgPath string `json:"img_path" gorm:"column:img_path;<-:false;->"`
	ImgName string `json:"img_name" gorm:"column:img_name;<-:false;->"`
}

func (Employee) TableName() string {
	return "m_employee"
}

// Based on EmployeeRepositoryInterface.php
type EmployeeRepository interface {
	FindAll(params dto.QueryParams, position string, majorId string) (*[]Employee, int64, error)
	FindByID(id string) (*Employee, error)
	FindAllAsOptions(position string, majorId string) (*[]Employee, error)
	Create(employee *Employee) (*Employee, error)
	Update(id string, employee *Employee) (*Employee, error)
	Delete(id string) error
}
