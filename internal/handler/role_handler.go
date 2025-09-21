package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	usecase usecase.RoleUseCase
}

func NewRoleHandler(usecase usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{usecase: usecase}
}

func (h *RoleHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	role, err := h.usecase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Role not found", err)
		return
	}

	roleResource := dto.RoleResource{
		ID:   role.ID,
		Name: role.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Role found", roleResource)
}

func (h *RoleHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	roles, totalRows, err := h.usecase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch roles", err)
		return
	}

	roleResources := []dto.RoleResource{}
	for _, role := range *roles {
		roleResources = append(roleResources, dto.RoleResource{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Roles fetched successfully", roleResources, meta)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var roleDTO dto.StoreRoleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	role, err := h.usecase.Create(&roleDTO)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			helper.ErrorResponse(c, http.StatusBadRequest, "Role name already exists", err)
			return
		}

		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create role", err)
		return
	}

	roleResource := dto.RoleResource{
		ID:   role.ID,
		Name: role.Name,
	}
	helper.SuccessResponse(c, http.StatusCreated, "Role created successfully", roleResource)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var roleDTO dto.UpdateRoleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	role, err := h.usecase.Update(id, &roleDTO)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Role not found", err)
			return
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			helper.ErrorResponse(c, http.StatusBadRequest, "Role name already exists", err)
			return
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update role", err)
			return
		}
	}

	roleResource := dto.RoleResource{
		ID:   role.ID,
		Name: role.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Role updated successfully", roleResource)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.Delete(id); err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Role not found", err)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete role", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Role deleted successfully", nil)
}
