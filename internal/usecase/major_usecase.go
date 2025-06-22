package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type MajorUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.Major, int64, error)
	FindByID(id string) (*domain.Major, error)
	FindAllAsOptions() (*[]domain.Major, error)
	Create(dto *dto.StoreMajorDTO) (*domain.Major, error)
	Update(id string, dto *dto.UpdateMajorDTO) (*domain.Major, error)
	Delete(id string) error
}

type majorUseCase struct {
	repo domain.MajorRepository
}

func NewMajorUseCase(repo domain.MajorRepository) MajorUseCase {
	return &majorUseCase{repo: repo}
}

func (u *majorUseCase) FindAll(params dto.QueryParams) (*[]domain.Major, int64, error) {
	return u.repo.FindAll(params)
}

func (u *majorUseCase) FindByID(id string) (*domain.Major, error) {
	return u.repo.FindByID(id)
}
func (u *majorUseCase) FindAllAsOptions() (*[]domain.Major, error) {
	return u.repo.FindAllAsOptions()
}

func (u *majorUseCase) Create(dto *dto.StoreMajorDTO) (*domain.Major, error) {
	major := &domain.Major{
		ID:   uuid.NewString(),
		Code: dto.Code,
		Name: dto.Name,
	}
	return u.repo.Create(major)
}

func (u *majorUseCase) Update(id string, dto *dto.UpdateMajorDTO) (*domain.Major, error) {
	major := &domain.Major{
		Code: dto.Code,
		Name: dto.Name,
	}
	return u.repo.Update(id, major)
}

func (u *majorUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
