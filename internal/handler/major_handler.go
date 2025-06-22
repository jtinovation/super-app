package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MajorHandler struct {
	useCase usecase.MajorUseCase
}

func NewMajorHandler(uc usecase.MajorUseCase) *MajorHandler {
	return &MajorHandler{useCase: uc}
}

func (h *MajorHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	majors, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch majors", err)
		return
	}

	majorResources := []dto.MajorResource{}
	for _, major := range *majors {
		majorResources = append(majorResources, dto.MajorResource{
			ID:   major.ID,
			Code: major.Code,
			Name: major.Name,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Majors fetched successfully", majorResources, meta)
}

func (h *MajorHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	major, err := h.useCase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Major not found", err)
		return
	}

	majorResource := dto.MajorResource{
		ID:   major.ID,
		Code: major.Code,
		Name: major.Name,
	}

	helper.SuccessResponse(c, http.StatusOK, "Major found", majorResource)
}

func (h *MajorHandler) Create(c *gin.Context) {
	var majorDTO dto.StoreMajorDTO
	if err := c.ShouldBindJSON(&majorDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	major, err := h.useCase.Create(&majorDTO)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create major", err)
		return
	}

	majorResource := dto.MajorResource{
		ID:   major.ID,
		Code: major.Code,
		Name: major.Name,
	}

	helper.SuccessResponse(c, http.StatusCreated, "Major created successfully", majorResource)
}

func (h *MajorHandler) FindAllAsOptions(c *gin.Context) {
	majors, err := h.useCase.FindAllAsOptions()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch majors", err)
		return
	}

	var optionResource []dto.Option
	for _, major := range *majors {
		optionResource = append(optionResource, dto.Option{
			Label: major.Name,
			Value: major.ID,
		})
	}

	helper.SuccessResponse(c, http.StatusOK, "Majors fetched successfully", optionResource)
}

func (h *MajorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var majorDTO dto.UpdateMajorDTO
	if err := c.ShouldBindJSON(&majorDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	major, err := h.useCase.Update(id, &majorDTO)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Major not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update major", err)
		}
		return
	}

	majorResource := dto.MajorResource{
		ID:   major.ID,
		Code: major.Code,
		Name: major.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Major updated successfully", majorResource)
}

func (h *MajorHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Major not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete major", err)
		}
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Major deleted successfully", nil)
}
