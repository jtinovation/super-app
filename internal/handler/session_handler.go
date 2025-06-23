package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	useCase usecase.SessionUseCase
}

func NewSessionHandler(uc usecase.SessionUseCase) *SessionHandler {
	return &SessionHandler{useCase: uc}
}

func (h *SessionHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	sessions, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch sessions", err)
		return
	}

	sessionResources := []dto.SessionResource{}
	for _, s := range *sessions {
		sessionResources = append(sessionResources, dto.SessionResource{
			ID:      s.ID,
			Session: s.Session,
		})
	}

	meta := &dto.Meta{Page: params.Page, PerPage: params.PerPage, Total: totalRows}
	helper.PaginatedSuccessResponse(c, http.StatusOK, "Sessions fetched successfully", sessionResources, meta)
}

func (h *SessionHandler) FindAllAsOptions(c *gin.Context) {
	sessions, err := h.useCase.FindAllAsOptions()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch session options", err)
		return
	}

	var optionResource []dto.Option
	for _, s := range *sessions {
		optionResource = append(optionResource, dto.Option{
			Label: s.Session,
			Value: s.ID,
		})
	}

	helper.SuccessResponse(c, http.StatusOK, "Session options fetched successfully", optionResource)
}

func (h *SessionHandler) Create(c *gin.Context) {
	var payload dto.StoreSessionDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	session, err := h.useCase.Create(&payload)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create session", err)
		return
	}

	resource := dto.SessionResource{ID: session.ID, Session: session.Session}
	helper.SuccessResponse(c, http.StatusCreated, "Session created successfully", resource)
}

func (h *SessionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var payload dto.UpdateSessionDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	session, err := h.useCase.Update(id, &payload)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Session not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update session", err)
		}
		return
	}

	resource := dto.SessionResource{ID: session.ID, Session: session.Session}
	helper.SuccessResponse(c, http.StatusOK, "Session updated successfully", resource)
}

func (h *SessionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Session not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete session", err)
		}
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Session deleted successfully", nil)
}
