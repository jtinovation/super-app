package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type EmployeeLab struct {
	ID         string  `gorm:"type:char(36);primaryKey"`
	LabID      string  `gorm:"column:m_lab_id;type:char(36);not null"`
	EmployeeID string  `gorm:"column:m_employee_id;type:char(36);not null"`
	IsHeadLab  bool    `gorm:"column:is_head_lab;type:boolean;not null"`
	Period     *string `gorm:"type:varchar(255);not null"`
	Status     string  `gorm:"type:enum('ACTIVE','INACTIVE');default:'ACTIVE'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	Lab      Lab      `gorm:"foreignKey:LabID;references:ID"`
	Employee Employee `gorm:"foreignKey:EmployeeID;references:ID"`
}

func (EmployeeLab) TableName() string {
	return "m_employee_lab"
}

type EmployeeLabRepository interface {
	FindAll(params dto.QueryParams, majorId string) (*[]EmployeeLab, int64, error)
	FindByID(id string) (*EmployeeLab, error)
	Create(employeeLab *EmployeeLab) error
	GetByID(id string) (*EmployeeLab, error)
	Update(employeeLab *EmployeeLab) error
	Delete(id string) error
}
