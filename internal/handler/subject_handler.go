package handler

import (
	"fmt"
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubjectHandler struct{ useCase usecase.SubjectUseCase }

func NewSubjectHandler(uc usecase.SubjectUseCase) *SubjectHandler {
	return &SubjectHandler{useCase: uc}
}

func (h *SubjectHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}
	subjects, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch subjects", err)
		return
	}
	resources := []dto.SubjectResource{}
	for _, s := range *subjects {
		resources = append(resources, dto.SubjectResource{
			ID: s.ID, Code: s.Code, Name: s.Name, Status: s.Status,
			StudyProgramName: s.StudyProgramName, StudyProgramID: s.StudyProgramID2,
		})
	}
	meta := &dto.Meta{Page: params.Page, PerPage: params.PerPage, Total: totalRows}
	helper.PaginatedSuccessResponse(c, http.StatusOK, "Subjects fetched successfully", resources, meta)
}

func (h *SubjectHandler) FindAllAsOptions(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")
	semesterID := c.Query("semester_id")
	subjects, err := h.useCase.FindAllAsOptions(studyProgramID, semesterID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch subject options", err)
		return
	}
	resources := []dto.Option{}
	for _, s := range *subjects {
		resources = append(resources, dto.Option{
			Value: s.ID, Label: fmt.Sprintf("(%s) %s", s.Code, s.Name),
		})
	}
	helper.SuccessResponse(c, http.StatusOK, "Subject options fetched successfully", resources)
}

func (h *SubjectHandler) Create(c *gin.Context) {
	var payload dto.StoreSubjectDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	subject, err := h.useCase.Create(&payload)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create subject", err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "Subject created successfully", subject)
}
func (h *SubjectHandler) Update(c *gin.Context) {
	var payload dto.UpdateSubjectDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	subjectID := c.Param("id")
	if subjectID == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "Subject ID is required", nil)
		return
	}

	subject, err := h.useCase.Update(subjectID, &payload)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update subject", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Subject updated successfully", subject)
}
func (h *SubjectHandler) Delete(c *gin.Context) {
	subjectID := c.Param("id")
	if subjectID == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "Subject ID is required", nil)
		return
	}

	err := h.useCase.Delete(subjectID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete subject", err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Subject deleted successfully", nil)
}

func (h *SubjectHandler) GetLectureOnSubject(c *gin.Context) {
	studyProgramID := c.Query("study_program_id")
	semesterID := c.Query("semester_id")
	if studyProgramID == "" || semesterID == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "study_program_id and semester_id are required", nil)
		return
	}

	data, err := h.useCase.GetLectureOnSubject(studyProgramID, semesterID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch lecture on subject data", err)
		return
	}

	resources := []dto.LectureOnSubjectResource{}
	for _, item := range *data {
		lectures := []dto.LectureResource{}
		for _, l := range item.Lecturers {
			avatarURL := ""
			if l.User.ImgPath != nil && l.User.ImgName != nil {
				avatarURL = helper.GetUrlFile(*l.User.ImgPath, *l.User.ImgName)
			}
			lectures = append(lectures, dto.LectureResource{
				ID: l.ID, MajorID: l.MajorID,
				User: dto.LectureUserResource{ID: l.User.ID, Name: l.User.Name, Avatar: avatarURL},
			})
		}
		resources = append(resources, dto.LectureOnSubjectResource{
			ID: item.ID, SemesterID: item.SemesterID,
			Subject: dto.LectureOnSubjectSubjectResource{
				ID: item.Subject.ID, Name: item.Subject.Name, Code: item.Subject.Code, StudyProgramID: item.Subject.StudyProgramID,
			},
			Lectures: lectures,
		})
	}
	helper.SuccessResponse(c, http.StatusOK, "Data fetched successfully", resources)
}

func (h *SubjectHandler) StoreLectureOnSubject(c *gin.Context) {
	var payload dto.SettingLectureOnSubjectDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}
	err := h.useCase.StoreLectureOnSubject(payload.Data)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to store lecture on subject", err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Lecturers assigned successfully", nil)
}
