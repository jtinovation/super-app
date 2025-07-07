package domain

import (
	"time"
)

type StudentSemester struct {
	ID         string `gorm:"type:char(36);primaryKey"`
	StudentID  string `gorm:"column:m_student_id;type:char(36);not null"`
	SemesterID string `gorm:"column:m_semester_id;type:char(36);not null"`
	Class      string `gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (StudentSemester) TableName() string {
	return "m_student_semester"
}

type StudentSemesterRepository interface {
	StoreStudentSemester(studentSemester *StudentSemester) error
}
