package dto

type RoleResource struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Permissions *[]PermissionResource `json:"permissions"`
}
type PermissionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StoreRoleDTO struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Permissions []string `json:"permissions" binding:"omitempty,dive,required"`
}

type UpdateRoleDTO struct {
	Name        string   `json:"name" binding:"required,max=100"`
	Permissions []string `json:"permissions" binding:"omitempty,dive,required"`
}

type StorePermissionDTO struct {
	Name string `json:"name" binding:"required,max=100"`
}

type UpdatePermissionDTO struct {
	Name string `json:"name" binding:"required,max=100"`
}

type RoleOptionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
