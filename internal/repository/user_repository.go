package repository

import (
	"jti-super-app-go/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) (*domain.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(id string, user *domain.User) (*domain.User, error) {
	var existingUser domain.User
	if err := r.db.First(&existingUser, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingUser).Updates(user).Error; err != nil {
		return nil, err
	}
	return &existingUser, nil
}
