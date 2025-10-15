package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubjectLectureHandler struct{ useCase usecase.SubjectLectureUseCase }

func NewSubjectLectureHandler(uc usecase.SubjectLectureUseCase) *SubjectLectureHandler {
	return &SubjectLectureHandler{useCase: uc}
}

func (h *SubjectLectureHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}

	result, totalRows, err := h.useCase.FindAll(*params)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch subject lectures", err)
		return
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Subject options fetched successfully", result, meta)
}
