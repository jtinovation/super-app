package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
)

type PermissionUseCase interface {
	FindByID(id string) (*domain.Permission, error)
	FindAll(params dto.QueryParams) (*[]domain.Permission, int64, error)
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
