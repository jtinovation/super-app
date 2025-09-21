package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"

	"github.com/google/uuid"
)

type RoleUseCase interface {
	FindByID(id string) (*domain.Role, error)
	FindAll(params dto.QueryParams) (*[]domain.Role, int64, error)
	Create(role *dto.StoreRoleDTO) (*domain.Role, error)
	Update(id string, role *dto.UpdateRoleDTO) (*domain.Role, error)
	Delete(id string) error
}

type roleUseCase struct {
	repo domain.RoleRepository
}

func NewRoleUseCase(repo domain.RoleRepository) RoleUseCase {
	return &roleUseCase{repo: repo}
}

func (u *roleUseCase) FindByID(id string) (*domain.Role, error) {
	return u.repo.FindByID(id)
}

func (u *roleUseCase) FindAll(params dto.QueryParams) (*[]domain.Role, int64, error) {
	return u.repo.FindAll(params)
}

func (u *roleUseCase) Create(dto *dto.StoreRoleDTO) (*domain.Role, error) {
	permissions := []domain.Permission{}
	for _, permID := range dto.Permissions {
		permissions = append(permissions, domain.Permission{ID: permID})
	}
	role := &domain.Role{
		ID:          uuid.NewString(),
		Name:        dto.Name,
		GuardName:   "web",
		Permissions: permissions,
	}

	return u.repo.Create(role)
}

func (u *roleUseCase) Update(id string, role *dto.UpdateRoleDTO) (*domain.Role, error) {
	permissions := []domain.Permission{}
	for _, permID := range role.Permissions {
		permissions = append(permissions, domain.Permission{ID: permID})
	}
	updatedRole := &domain.Role{
		Name:        role.Name,
		Permissions: permissions,
	}
	return u.repo.Update(id, updatedRole)
}

func (u *roleUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
