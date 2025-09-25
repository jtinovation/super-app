package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(uc usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: uc}
}

func (h *UserHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	users, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", err)
		return
	}

	userResources := []dto.UserResource{}
	for _, user := range *users {
		roles := make([]dto.RoleOptionResource, len(user.Roles))
		for i, role := range user.Roles {
			roles[i] = dto.RoleOptionResource{
				ID:   role.ID,
				Name: role.Name,
			}
		}
		userResources = append(userResources, dto.UserResource{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Status: user.Status,
			Roles:  roles,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Users fetched successfully", userResources, meta)
}

func (h *UserHandler) UpdateRoles(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := h.useCase.UpdateRoles(id, req); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user roles", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "User roles updated successfully", nil)
}
