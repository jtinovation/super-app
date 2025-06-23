package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SemesterHandler struct {
	useCase usecase.SemesterUseCase
}

func NewSemesterHandler(uc usecase.SemesterUseCase) *SemesterHandler {
	return &SemesterHandler{useCase: uc}
}

func (h *SemesterHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	semesters, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch semesters", err)
		return
	}

	semesterResources := []dto.SemesterResource{}
	for _, s := range *semesters {
		semesterResources = append(semesterResources, dto.SemesterResource{
			ID:       s.ID,
			Year:     s.Year,
			Semester: s.Semester,
			Session:  dto.SessionResource{ID: s.SessionID, Session: s.SessionName},
		})
	}
	meta := &dto.Meta{Page: params.Page, PerPage: params.PerPage, Total: totalRows}
	helper.PaginatedSuccessResponse(c, http.StatusOK, "Semesters fetched successfully", semesterResources, meta)
}

func (h *SemesterHandler) FindAllAsOptions(c *gin.Context) {
	sessionID := c.Query("session_id")
	semesters, err := h.useCase.FindAllAsOptions(sessionID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch semester options", err)
		return
	}

	optionResource := []dto.SemesterOptionResource{}
	for _, s := range *semesters {
		optionResource = append(optionResource, dto.SemesterOptionResource{
			ID:       s.ID,
			Year:     s.Year,
			Semester: s.Semester,
		})
	}
	helper.SuccessResponse(c, http.StatusOK, "Semester options fetched successfully", optionResource)
}

func (h *SemesterHandler) Create(c *gin.Context) {
	var payload dto.StoreSemesterDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}
	semester, err := h.useCase.Create(&payload)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create semester", err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "Semester created successfully", semester)
}

func (h *SemesterHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var payload dto.UpdateSemesterDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}
	semester, err := h.useCase.Update(id, &payload)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update semester", err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Semester updated successfully", semester)
}

func (h *SemesterHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.useCase.Delete(id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete semester", err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Semester deleted successfully", nil)
}

func (h *SemesterHandler) SettingSubjectSemester(c *gin.Context) {
	id := c.Param("id")
	var payload dto.SettingSubjectSemesterDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}
	if err := h.useCase.SettingSubjectSemester(id, payload.SubjectIDs); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to set subjects for semester", err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Subjects for semester updated successfully", nil)
}
