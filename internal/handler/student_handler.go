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
			UserID:       s.UserID,
			NIM:          s.NIM,
			Name:         s.Name,
			Generation:   s.Generation,
			Avatar:       avatarURL,
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
		maxSize := int64(2 * 1024 * 1024) // 2MB
		allowedMimeTypes := map[string]bool{
			"image/jpeg":    true,
			"image/png":     true,
			"image/gif":     true,
			"image/svg+xml": true,
		}

		if !helper.ValidateUploadedFile(c, file, maxSize, allowedMimeTypes) {
			return
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

func (h *StudentHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "Student ID is required", nil)
		return
	}

	student, err := h.useCase.FindByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			helper.ErrorResponse(c, http.StatusNotFound, "Student not found", nil)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch student", err)
		return
	}

	avatarURL := ""
	if student.ImgPath != "" && student.ImgName != "" {
		avatarURL = helper.GetUrlFile(student.ImgPath, student.ImgName)
	}

	semesters := []dto.StudentSemesterResource{}
	for _, ss := range student.StudentSemesters {
		semesters = append(semesters, dto.StudentSemesterResource{
			ID:        ss.Semester.ID,
			Year:      ss.Semester.Year,
			Semester:  ss.Semester.Semester,
			Class:     ss.Class,
			SessionId: ss.Semester.SessionID,
			Session:   ss.Semester.Session.Session,
		})
	}

	resource := dto.StudentDetailResource{
		ID:             student.ID,
		NIM:            student.NIM,
		Generation:     student.Generation,
		TuitionFee:     student.TuitionFee,
		TuitionMethod:  student.TuitionMethod,
		StudyProgramId: student.StudyProgram.ID,
		MajorId:        student.StudyProgram.MajorID,
		User: dto.UserResource{
			ID:          student.User.ID,
			Name:        student.User.Name,
			Email:       student.User.Email,
			Status:      student.User.Status,
			Gender:      student.User.Gender,
			Religion:    student.User.Religion,
			BirthDate:   student.User.BirthDate,
			BirthPlace:  student.User.BirthPlace,
			PhoneNumber: student.User.PhoneNumber,
			Address:     student.User.Address,
			Nationality: student.User.Nationality,
			Avatar:      avatarURL,
		},
		Semesters: semesters,
	}

	helper.SuccessResponse(c, http.StatusOK, "Student fetched successfully", resource)
}
