package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OauthClientHandler struct {
	useCase usecase.OauthClientUseCase
}

func NewOauthClientHandler(uc usecase.OauthClientUseCase) *OauthClientHandler {
	return &OauthClientHandler{useCase: uc}
}

func (h *OauthClientHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	clients, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch OAuth clients", err)
		return
	}

	clientResources := []dto.OauthClientResource{}
	for _, client := range *clients {
		clientResources = append(clientResources, dto.OauthClientResource{
			ID:       client.ID,
			Name:     client.Name,
			Redirect: client.Redirect,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "OAuth clients fetched successfully", clientResources, meta)
}

func (h *OauthClientHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	client, err := h.useCase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch OAuth client", err)
		return
	}

	clientResource := dto.OauthClientResource{
		ID:       client.ID,
		Name:     client.Name,
		Redirect: client.Redirect,
	}

	helper.SuccessResponse(c, http.StatusOK, "OAuth client fetched successfully", clientResource)
}

func (h *OauthClientHandler) Create(c *gin.Context) {
	var clientDTO dto.StoreOauthClientDTO
	if err := c.ShouldBindJSON(&clientDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	client, err := h.useCase.Create(&clientDTO)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create OAuth client", err)
		return
	}

	clientResource := dto.OauthClientResource{
		ID:       client.ID,
		Name:     client.Name,
		Redirect: client.Redirect,
	}

	helper.SuccessResponse(c, http.StatusCreated, "OAuth client created successfully", clientResource)
}

func (h *OauthClientHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var clientDTO dto.UpdateOauthClientDTO
	if err := c.ShouldBindJSON(&clientDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	client, err := h.useCase.Update(id, &clientDTO)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update OAuth client", err)
		return
	}

	clientResource := dto.OauthClientResource{
		ID:       client.ID,
		Name:     client.Name,
		Redirect: client.Redirect,
	}

	helper.SuccessResponse(c, http.StatusOK, "OAuth client updated successfully", clientResource)
}

func (h *OauthClientHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete OAuth client", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "OAuth client deleted successfully", nil)
}
