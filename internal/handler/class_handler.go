package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	useCase usecase.ClassUseCase
}

func NewClassHandler(uc usecase.ClassUseCase) *ClassHandler {
	return &ClassHandler{useCase: uc}
}

func (h *ClassHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	majorId := c.Query("major_id")
	studyProgramId := c.Query("study_program_id")
	if !ok {
		return
	}

	classs, totalRows, err := h.useCase.FindAll(*params, studyProgramId, majorId)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch classs", err)
		return
	}

	classResources := []dto.ClassResource{}
	for _, class := range *classs {
		classResources = append(classResources, dto.ClassResource{
			ID:   class.ID,
			Code: class.Code,
			Name: class.Name,
			StudyProgram: dto.StudyProgramOptionResource{
				ID:   class.StudyProgramID,
				Name: class.StudyProgramName,
			},
			Major: dto.MajorOptionResource{
				ID:   class.MajorID,
				Name: class.MajorName,
			},
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Classs fetched successfully", classResources, meta)
}

func (h *ClassHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	class, err := h.useCase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Class not found", err)
		return
	}

	classResource := dto.ClassResource{
		ID:   class.ID,
		Code: class.Code,
		Name: class.Name,
	}

	helper.SuccessResponse(c, http.StatusOK, "Class found", classResource)
}

func (h *ClassHandler) Create(c *gin.Context) {
	var classDTO dto.StoreClassDTO
	if err := c.ShouldBindJSON(&classDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	class, err := h.useCase.Create(&classDTO)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create class", err)
		return
	}

	classResource := dto.ClassResource{
		ID:   class.ID,
		Code: class.Code,
		Name: class.Name,
	}

	helper.SuccessResponse(c, http.StatusCreated, "Class created successfully", classResource)
}

func (h *ClassHandler) FindAllAsOptions(c *gin.Context) {
	classs, err := h.useCase.FindAllAsOptions()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch classs", err)
		return
	}

	var optionResource []dto.Option
	for _, class := range *classs {
		optionResource = append(optionResource, dto.Option{
			Label: class.Name,
			Value: class.ID,
		})
	}

	helper.SuccessResponse(c, http.StatusOK, "Classs fetched successfully", optionResource)
}

func (h *ClassHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var classDTO dto.UpdateClassDTO
	if err := c.ShouldBindJSON(&classDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	class, err := h.useCase.Update(id, &classDTO)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Class not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update class", err)
		}
		return
	}

	classResource := dto.ClassResource{
		ID:   class.ID,
		Code: class.Code,
		Name: class.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Class updated successfully", classResource)
}

func (h *ClassHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Class not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete class", err)
		}
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Class deleted successfully", nil)
}
