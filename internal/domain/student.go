package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID               string  `gorm:"type:char(36);primaryKey"`
	UserID           string  `gorm:"column:m_user_id;type:char(36);not null"`
	StudentProgramID string  `gorm:"column:m_study_program_id;type:char(36);not null"`
	NIM              string  `gorm:"type:varchar(255);not null;unique"`
	Generation       *int    `gorm:"type:int"`
	TuitionFee       *int    `gorm:"type:int"`
	TuitionMethod    *string `gorm:"type:varchar(255)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	User             User              `gorm:"foreignKey:UserID;references:ID"`
	StudyProgram     StudyProgram      `gorm:"foreignKey:StudentProgramID;references:ID"`
	StudentSemesters []StudentSemester `gorm:"foreignKey:StudentID;references:ID"`
	// Semesters       []Semester        `gorm:"many2many:m_student_semester;joinForeignKey:m_student_id;joinReferences:m_semester_id;References:ID;foreignKey:ID"`

	Name             string `gorm:"column:name;<-:false;->"`
	ImgPath          string `gorm:"column:img_path;<-:false;->"`
	ImgName          string `gorm:"column:img_name;<-:false;->"`
	StudyProgramName string `gorm:"column:study_program_name;<-:false;->"`
	StudyProgramID   string `gorm:"column:study_program_id;<-:false;->"`
	MajorName        string `gorm:"column:major_name;<-:false;->"`
	MajorID          string `gorm:"column:major_id;<-:false;->"`
}

func (Student) TableName() string {
	return "m_student"
}

type StudentRepository interface {
	FindAll(params dto.QueryParams) (*[]Student, int64, error)
	FindByID(id string) (*Student, error)
	Create(student *Student) (*Student, error)
	Update(id string, student *Student) (*Student, error)
	Delete(id string) error
}
