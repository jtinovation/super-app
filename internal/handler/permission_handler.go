package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	usecase usecase.PermissionUseCase
}

func NewPermissionHandler(usecase usecase.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{usecase: usecase}
}

func (h *PermissionHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	permission, err := h.usecase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Permission not found", err)
		return
	}

	permissionResource := dto.PermissionResource{
		ID:   permission.ID,
		Name: permission.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Permission found", permissionResource)
}

func (h *PermissionHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	permissions, totalRows, err := h.usecase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permissions", err)
		return
	}

	permissionResources := []dto.PermissionResource{}
	for _, permission := range *permissions {
		permissionResources = append(permissionResources, dto.PermissionResource{
			ID:   permission.ID,
			Name: permission.Name,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Permissions fetched successfully", permissionResources, meta)
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var input dto.StorePermissionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	permission, err := h.usecase.Create(&input)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			helper.ErrorResponse(c, http.StatusBadRequest, "Permission with this name already exists", err)
			return
		}

		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create permission", err)
		return
	}

	permissionResource := dto.PermissionResource{
		ID:   permission.ID,
		Name: permission.Name,
	}
	helper.SuccessResponse(c, http.StatusCreated, "Permission created successfully", permissionResource)
}

func (h *PermissionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.UpdatePermissionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	permission, err := h.usecase.Update(id, &input)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Permission not found", err)
			return
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			helper.ErrorResponse(c, http.StatusBadRequest, "Permission with this name already exists", err)
			return
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update permission", err)
			return
		}
	}

	permissionResource := dto.PermissionResource{
		ID:   permission.ID,
		Name: permission.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Permission updated successfully", permissionResource)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.Delete(id); err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Permission not found", err)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete permission", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Permission deleted successfully", nil)
}
