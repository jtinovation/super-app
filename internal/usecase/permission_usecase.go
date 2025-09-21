package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type PermissionUseCase interface {
	FindByID(id string) (*domain.Permission, error)
	FindAll(params dto.QueryParams) (*[]domain.Permission, int64, error)
	Create(dto *dto.StorePermissionDTO) (*domain.Permission, error)
	Update(id string, permission *dto.UpdatePermissionDTO) (*domain.Permission, error)
	Delete(id string) error
}

type permissionUseCase struct {
	repo domain.PermissionRepository
}

func NewPermissionUseCase(repo domain.PermissionRepository) PermissionUseCase {
	return &permissionUseCase{repo: repo}
}

func (u *permissionUseCase) FindByID(id string) (*domain.Permission, error) {
	return u.repo.FindByID(id)
}

func (u *permissionUseCase) FindAll(params dto.QueryParams) (*[]domain.Permission, int64, error) {
	return u.repo.FindAll(params)
}

func (u *permissionUseCase) Create(dto *dto.StorePermissionDTO) (*domain.Permission, error) {
	permission := &domain.Permission{
		ID:        uuid.NewString(),
		Name:      dto.Name,
		GuardName: "web",
	}

	return u.repo.Create(permission)
}

func (u *permissionUseCase) Update(id string, dto *dto.UpdatePermissionDTO) (*domain.Permission, error) {
	permission := &domain.Permission{
		Name: dto.Name,
	}

	return u.repo.Update(id, permission)
}

func (u *permissionUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
