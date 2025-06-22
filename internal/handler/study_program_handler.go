package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StudyProgramHandler struct {
	useCase usecase.StudyProgramUseCase
}

func NewStudyProgramHandler(uc usecase.StudyProgramUseCase) *StudyProgramHandler {
	return &StudyProgramHandler{useCase: uc}
}

func (h *StudyProgramHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	majorId := c.Query("major_id")
	if !ok {
		return
	}

	studyPrograms, totalRows, err := h.useCase.FindAll(*params, majorId)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch studyPrograms", err)
		return
	}

	studyProgramResources := []dto.StudyProgramResource{}
	for _, studyProgram := range *studyPrograms {
		studyProgramResources = append(studyProgramResources, dto.StudyProgramResource{
			ID:   studyProgram.ID,
			Code: studyProgram.Code,
			Name: studyProgram.Name,
			Major: dto.MajorOptionResource{
				ID:   studyProgram.MajorID,
				Name: studyProgram.MajorName,
			},
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "StudyPrograms fetched successfully", studyProgramResources, meta)
}

func (h *StudyProgramHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	studyProgram, err := h.useCase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "StudyProgram not found", err)
		return
	}

	studyProgramResource := dto.StudyProgramResource{
		ID:   studyProgram.ID,
		Code: studyProgram.Code,
		Name: studyProgram.Name,
	}

	helper.SuccessResponse(c, http.StatusOK, "StudyProgram found", studyProgramResource)
}

func (h *StudyProgramHandler) Create(c *gin.Context) {
	var studyProgramDTO dto.StoreStudyProgramDTO
	if err := c.ShouldBindJSON(&studyProgramDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	studyProgram, err := h.useCase.Create(&studyProgramDTO)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create studyProgram", err)
		return
	}

	studyProgramResource := dto.StudyProgramResource{
		ID:   studyProgram.ID,
		Code: studyProgram.Code,
		Name: studyProgram.Name,
	}

	helper.SuccessResponse(c, http.StatusCreated, "StudyProgram created successfully", studyProgramResource)
}

func (h *StudyProgramHandler) FindAllAsOptions(c *gin.Context) {
	studyPrograms, err := h.useCase.FindAllAsOptions()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch studyPrograms", err)
		return
	}

	var optionResource []dto.Option
	for _, studyProgram := range *studyPrograms {
		optionResource = append(optionResource, dto.Option{
			Label: studyProgram.Name,
			Value: studyProgram.ID,
		})
	}

	helper.SuccessResponse(c, http.StatusOK, "StudyPrograms fetched successfully", optionResource)
}

func (h *StudyProgramHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var studyProgramDTO dto.UpdateStudyProgramDTO
	if err := c.ShouldBindJSON(&studyProgramDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	studyProgram, err := h.useCase.Update(id, &studyProgramDTO)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "StudyProgram not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update studyProgram", err)
		}
		return
	}

	studyProgramResource := dto.StudyProgramResource{
		ID:   studyProgram.ID,
		Code: studyProgram.Code,
		Name: studyProgram.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "StudyProgram updated successfully", studyProgramResource)
}

func (h *StudyProgramHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "StudyProgram not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete studyProgram", err)
		}
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "StudyProgram deleted successfully", nil)
}
