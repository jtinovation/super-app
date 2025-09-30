package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectLecture struct {
	ID                string   `gorm:"type:char(36);primaryKey"`
	SubjectSemesterID string   `gorm:"column:m_subject_semester_id;type:char(36);not null"`
	EmployeeID        string   `gorm:"column:m_employee_id;type:char(36);not null"`
	Employee          Employee `gorm:"foreignKey:EmployeeID"`
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
