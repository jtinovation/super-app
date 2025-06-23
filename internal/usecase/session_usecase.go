package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type SessionUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.Session, int64, error)
	FindByID(id string) (*domain.Session, error)
	FindAllAsOptions() (*[]domain.Session, error)
	Create(dto *dto.StoreSessionDTO) (*domain.Session, error)
	Update(id string, dto *dto.UpdateSessionDTO) (*domain.Session, error)
	Delete(id string) error
}

type sessionUseCase struct {
	repo domain.SessionRepository
}

func NewSessionUseCase(repo domain.SessionRepository) SessionUseCase {
	return &sessionUseCase{repo: repo}
}

func (u *sessionUseCase) FindAll(params dto.QueryParams) (*[]domain.Session, int64, error) {
	return u.repo.FindAll(params)
}

func (u *sessionUseCase) FindByID(id string) (*domain.Session, error) {
	return u.repo.FindByID(id)
}
func (u *sessionUseCase) FindAllAsOptions() (*[]domain.Session, error) {
	return u.repo.FindAllAsOptions()
}

func (u *sessionUseCase) Create(dto *dto.StoreSessionDTO) (*domain.Session, error) {
	session := &domain.Session{
		ID:      uuid.NewString(),
		Session: dto.Session,
	}
	return u.repo.Create(session)
}

func (u *sessionUseCase) Update(id string, dto *dto.UpdateSessionDTO) (*domain.Session, error) {
	session := &domain.Session{
		Session: dto.Session,
	}
	return u.repo.Update(id, session)
}

func (u *sessionUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
