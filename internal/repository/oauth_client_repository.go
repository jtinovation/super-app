package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type oauthClientRepository struct {
	db *gorm.DB
}

func NewOauthClientRepository(db *gorm.DB) domain.OauthClientRepository {
	return &oauthClientRepository{db: db}
}

func (r *oauthClientRepository) FindAll(params dto.QueryParams) (*[]domain.OauthClient, int64, error) {
	var clients []domain.OauthClient
	var totalRows int64

	query := r.db.Model(&domain.OauthClient{})

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where("LOWER(name) LIKE ?", searchQuery)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("name asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return &clients, totalRows, nil
}

func (r *oauthClientRepository) FindByID(id string) (*domain.OauthClient, error) {
	var client domain.OauthClient
	if err := r.db.Where("id = ?", id).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *oauthClientRepository) Create(client *domain.OauthClient) error {
	return r.db.Create(client).Error
}

func (r *oauthClientRepository) Update(id string, client *domain.OauthClient) error {
	return r.db.Model(&domain.OauthClient{}).Where("id = ?", id).Updates(client).Error
}

func (r *oauthClientRepository) Delete(id string) error {
	return r.db.Delete(&domain.OauthClient{}, "id = ?", id).Error
}
