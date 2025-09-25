package domain

import "jti-super-app-go/internal/dto"

type AuthRepository interface {
	FindByEmail(email string) (*User, error)
}

type UserRepository interface {
	FindAll(params dto.QueryParams) (*[]User, int64, error)
	FindByID(id string) (*User, error)
	Create(user *User) (*User, error)
	Update(id string, user *User) (*User, error)
	UpdateRoles(id string, roles []Role) error
}
