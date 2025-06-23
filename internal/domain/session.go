package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	Session   string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Session) TableName() string {
	return "m_session"
}

type SessionRepository interface {
	FindAll(params dto.QueryParams) (*[]Session, int64, error)
	FindByID(id string) (*Session, error)
	FindAllAsOptions() (*[]Session, error)
	Create(session *Session) (*Session, error)
	Update(id string, session *Session) (*Session, error)
	Delete(id string) error
}
