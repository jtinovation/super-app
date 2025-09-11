package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LabHandler struct {
	useCase usecase.LabUseCase
}

func NewLabHandler(uc usecase.LabUseCase) *LabHandler {
	return &LabHandler{useCase: uc}
}

func (h *LabHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	majorId := c.Query("major_id")
	if !ok {
		return
	}

	labs, totalRows, err := h.useCase.FindAll(*params, majorId)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch labs", err)
		return
	}

	labResources := []dto.LabResource{}
	for _, lab := range *labs {
		employeeLabResources := []dto.EmployeeLabResourceMany{}
		for _, empLab := range lab.EmployeeLab {
			employeeLabResources = append(employeeLabResources, dto.EmployeeLabResourceMany{
				ID:        empLab.ID,
				Name:      empLab.Employee.Name,
				IsHeadLab: empLab.IsHeadLab,
				Status:    empLab.Status,
			})
		}
		labResources = append(labResources, dto.LabResource{
			ID:   lab.ID,
			Code: lab.Code,
			Name: lab.Name,
			Major: dto.MajorOptionResource{
				ID:   lab.MajorID,
				Name: lab.MajorName,
			},
			EmployeeLab: employeeLabResources,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Labs fetched successfully", labResources, meta)
}

func (h *LabHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	lab, err := h.useCase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Lab not found", err)
		return
	}

	labResource := dto.LabResource{
		ID:   lab.ID,
		Code: lab.Code,
		Name: lab.Name,
	}

	helper.SuccessResponse(c, http.StatusOK, "Lab found", labResource)
}

func (h *LabHandler) Create(c *gin.Context) {
	var labDTO dto.StoreLabDTO
	if err := c.ShouldBindJSON(&labDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	lab, err := h.useCase.Create(&labDTO)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create lab", err)
		return
	}

	labResource := dto.LabResource{
		ID:   lab.ID,
		Code: lab.Code,
		Name: lab.Name,
		Major: dto.MajorOptionResource{
			ID:   lab.MajorID,
			Name: lab.MajorName,
		},
	}

	helper.SuccessResponse(c, http.StatusCreated, "Lab created successfully", labResource)
}

func (h *LabHandler) FindAllAsOptions(c *gin.Context) {
	majorId := c.Query("major_id")
	labs, err := h.useCase.FindAllAsOptions(majorId)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch labs", err)
		return
	}

	var optionResource []dto.Option
	for _, lab := range *labs {
		optionResource = append(optionResource, dto.Option{
			Label: lab.Name,
			Value: lab.ID,
		})
	}

	helper.SuccessResponse(c, http.StatusOK, "Labs fetched successfully", optionResource)
}

func (h *LabHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var labDTO dto.UpdateLabDTO
	if err := c.ShouldBindJSON(&labDTO); err != nil {
		helper.ValidationErrorJSON(c, err)
		return
	}

	lab, err := h.useCase.Update(id, &labDTO)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Lab not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update lab", err)
		}
		return
	}

	labResource := dto.LabResource{
		ID:   lab.ID,
		Code: lab.Code,
		Name: lab.Name,
	}
	helper.SuccessResponse(c, http.StatusOK, "Lab updated successfully", labResource)
}

func (h *LabHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Lab not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete lab", err)
		}
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Lab deleted successfully", nil)
}
