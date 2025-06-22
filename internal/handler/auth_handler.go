package handler

import (
	"errors"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: uc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	res, err := h.useCase.Login(req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Could not retrieve token from context", errors.New("token not found in context"))
		return
	}

	tokenString, ok := token.(string)
	if !ok {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Invalid token format in context", errors.New("invalid token format"))
		return
	}

	err := h.useCase.Logout(tokenString)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error(), errors.New("logout failed"))
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}