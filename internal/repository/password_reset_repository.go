package repository

import (
	"jti-super-app-go/internal/domain"

	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) domain.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(pr *domain.PasswordReset) (*domain.PasswordReset, error) {
	r.db.Where("email = ?", pr.Email).Delete(&domain.PasswordReset{})

	if err := r.db.Create(pr).Error; err != nil {
		return nil, err
	}
	return pr, nil
}

func (r *passwordResetRepository) FindByTokenAndEmail(token string, email string) (*domain.PasswordReset, error) {
	var pr domain.PasswordReset
	if err := r.db.Where("token = ? AND email = ?", token, email).First(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *passwordResetRepository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.PasswordReset{}).Error
}
