package domain

import "jti-super-app-go/internal/dto"

type RoleRepository interface {
	FindByID(id string) (*Role, error)
	FindAll(params dto.QueryParams) (*[]Role, int64, error)
	Create(role *Role) (*Role, error)
	Update(id string, role *Role) (*Role, error)
	Delete(id string) error
}

type PermissionRepository interface {
	FindByID(id string) (*Permission, error)
	FindAll(params dto.QueryParams) (*[]Permission, int64, error)
	Create(permission *Permission) (*Permission, error)
	Update(id string, permission *Permission) (*Permission, error)
	Delete(id string) error
}
