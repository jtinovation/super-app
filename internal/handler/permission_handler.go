package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

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
