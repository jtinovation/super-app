package domain

import (
	"time"
)

type PasswordReset struct {
	Email     string `gorm:"primaryKey"`
	Token     string `gorm:"index"`
	CreatedAt time.Time
}

func (PasswordReset) TableName() string {
	return "password_reset_tokens"
}

type PasswordResetRepository interface {
	Create(pr *PasswordReset) (*PasswordReset, error)
	FindByTokenAndEmail(token string, email string) (*PasswordReset, error)
	Delete(token string) error
}
