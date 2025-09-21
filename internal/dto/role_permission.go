package dto

type RoleResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type PermissionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StoreRoleDTO struct {
	Name string `json:"name" binding:"required,max=100"`
}

type UpdateRoleDTO struct {
	Name string `json:"name" binding:"required,max=100"`
}

type StorePermissionDTO struct {
	Name string `json:"name" binding:"required,max=100"`
}

type UpdatePermissionDTO struct {
	Name string `json:"name" binding:"required,max=100"`
}
