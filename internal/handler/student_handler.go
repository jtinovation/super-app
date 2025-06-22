package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	useCase usecase.StudentUseCase
}

func NewStudentHandler(uc usecase.StudentUseCase) *StudentHandler {
	return &StudentHandler{useCase: uc}
}

func (h *StudentHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	students, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch students", err)
		return
	}

	studentResources := []dto.StudentResource{}
	for _, s := range *students {
		avatarURL := ""
		if s.ImgPath != "" && s.ImgName != "" {
			avatarURL = helper.GetUrlFile(s.ImgPath, s.ImgName)
		}
		studentResources = append(studentResources, dto.StudentResource{
			ID:           s.ID,
			NIM:          s.NIM,
			Name:         s.Name,
			Generation:   s.Generation,
			Avatar:       avatarURL,
			Class:        dto.ClassOptionResource{ID: s.ClassID, Name: s.ClassName},
			StudyProgram: dto.StudyProgramOptionResource{ID: s.StudyProgramID, Name: s.StudyProgramName},
			Major:        dto.MajorOptionResource{ID: s.MajorID, Name: s.MajorName},
		})
	}

	meta := &dto.Meta{Page: params.Page, PerPage: params.PerPage, Total: totalRows}
	helper.PaginatedSuccessResponse(c, http.StatusOK, "Students fetched successfully", studentResources, meta)
}

func (h *StudentHandler) Create(c *gin.Context) {
	var payload dto.StoreStudentDTO
	if err := c.ShouldBind(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	file, err := c.FormFile("avatar")
	if err == nil {
		if !helper.ValidateUploadedFile(c, file, 2*1024*1024, map[string]bool{
			"image/jpeg": true, "image/png": true, "image/gif": true, "image/svg+xml": true,
		}) {
			return // Validation failed
		}
		payload.Avatar = file
	} else if err != http.ErrMissingFile {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid file upload", err)
		return
	}

	student, err := h.useCase.Create(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			helper.ErrorResponse(c, http.StatusConflict, "Student with this NIM already exists", err)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create student", err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "Student created successfully", student)
}
