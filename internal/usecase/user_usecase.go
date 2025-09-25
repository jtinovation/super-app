package usecase

import (
	"jti-super-app-go/internal/domain"
	"jti-super-app-go/internal/dto"
)

type UserUseCase interface {
	FindAll(params dto.QueryParams) (*[]domain.User, int64, error)
	UpdateRoles(id string, roles dto.UpdateUserRolesRequest) error
}

type userUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) FindAll(params dto.QueryParams) (*[]domain.User, int64, error) {
	return u.repo.FindAll(params)
}

func (u *userUseCase) UpdateRoles(id string, roles dto.UpdateUserRolesRequest) error {
	roleList := make([]domain.Role, len(roles.RoleIDs))
	for i, roleID := range roles.RoleIDs {
		roleList[i] = domain.Role{ID: roleID}
	}

	return u.repo.UpdateRoles(id, roleList)
}
