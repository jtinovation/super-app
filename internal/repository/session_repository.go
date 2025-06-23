package repository

import (
	"fmt"
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
	"strings"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) FindAll(params dto.QueryParams) (*[]domain.Session, int64, error) {
	var sessions []domain.Session
	var totalRows int64

	query := r.db.Model(&domain.Session{})

	if params.Search != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(params.Search))
		query = query.Where("LOWER(session) LIKE ?", searchQuery)
	}

	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	if params.Sort != "" {
		sortOrder := fmt.Sprintf("%s %s", params.Sort, params.Order)
		query = query.Order(sortOrder)
	} else {
		query = query.Order("session asc")
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Offset(offset).Limit(params.PerPage)

	if err := query.Find(&sessions).Error; err != nil {
		return nil, 0, err
	}

	return &sessions, totalRows, nil
}

func (r *sessionRepository) FindByID(id string) (*domain.Session, error) {
	var session domain.Session
	if err := r.db.First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) FindAllAsOptions() (*[]domain.Session, error) {
	var sessions []domain.Session
	// In handler, we will map this to dto.Option
	if err := r.db.Select("id", "session").Order("session asc").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return &sessions, nil
}

func (r *sessionRepository) Create(session *domain.Session) (*domain.Session, error) {
	if err := r.db.Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (r *sessionRepository) Update(id string, session *domain.Session) (*domain.Session, error) {
	var existingSession domain.Session
	if err := r.db.First(&existingSession, "id = ?", id).Error; err != nil {
		return nil, err // record not found
	}

	if err := r.db.Model(&existingSession).Updates(session).Error; err != nil {
		return nil, err
	}
	return &existingSession, nil
}

func (r *sessionRepository) Delete(id string) error {
	if err := r.db.First(&domain.Session{}, "id = ?", id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Session{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
