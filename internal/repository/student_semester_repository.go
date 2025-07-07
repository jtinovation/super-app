package repository

import (
	"jti-super-app-go/internal/domain"

	"gorm.io/gorm"
)

type studentSemesterRepository struct {
	db *gorm.DB
}

func NewStudentSemesterRepository(db *gorm.DB) domain.StudentSemesterRepository {
	return &studentSemesterRepository{
		db: db,
	}
}

func (r *studentSemesterRepository) StoreStudentSemester(studentSemester *domain.StudentSemester) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(studentSemester).Error; err != nil {
			return err
		}

		return nil
	})
}
