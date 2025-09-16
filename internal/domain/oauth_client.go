package domain

import (
	"jti-super-app-go/internal/dto"
	"time"

	"gorm.io/gorm"
)

type OauthClient struct {
	ID        string `gorm:"type:char(36);primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	Secret    string `gorm:"type:varchar(100);not null"`
	Redirect  string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (OauthClient) TableName() string {
	return "m_oauth_client"
}

type OauthClientRepository interface {
	FindAll(params dto.QueryParams) (*[]OauthClient, int64, error)
	FindByID(id string) (*OauthClient, error)
	Create(client *OauthClient) error
	Update(id string, client *OauthClient) error
	Delete(id string) error
}
