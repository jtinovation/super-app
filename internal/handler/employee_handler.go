package handler

import (
	"jti-super-app-go/internal/dto"
	"jti-super-app-go/internal/usecase"
	"jti-super-app-go/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	useCase usecase.EmployeeUseCase
}

func NewEmployeeHandler(uc usecase.EmployeeUseCase) *EmployeeHandler {
	return &EmployeeHandler{useCase: uc}
}

func (h *EmployeeHandler) FindAll(c *gin.Context) {
	params, ok := helper.ParsePaginationQuery(c)
	if !ok {
		return
	}
	position := c.Query("position")
	majorId := c.Query("major_id")

	employees, totalRows, err := h.useCase.FindAll(*params, position, majorId)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch employees", err)
		return
	}

	employeeResources := []dto.EmployeeResource{}
	for _, emp := range *employees {
		avatarURL := ""
		if emp.ImgPath != "" && emp.ImgName != "" {
			avatarURL = helper.GetUrlFile(emp.ImgPath, emp.ImgName)
		}

		employeeResources = append(employeeResources, dto.EmployeeResource{
			ID:       emp.ID,
			UserID:   emp.UserID,
			Name:     emp.Name,
			Email:    emp.Email,
			NIP:      emp.Nip,
			Position: emp.Position,
			Avatar:   avatarURL,
		})
	}

	meta := &dto.Meta{
		Page:    params.Page,
		PerPage: params.PerPage,
		Total:   totalRows,
	}

	helper.PaginatedSuccessResponse(c, http.StatusOK, "Employees fetched successfully", employeeResources, meta)
}

func (h *EmployeeHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	employee, err := h.useCase.FindByID(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Employee not found", err)
		return
	}

	avatarURL := ""
	if employee.User.ImgPath != nil && employee.User.ImgName != nil {
		avatarURL = helper.GetUrlFile(*employee.User.ImgPath, *employee.User.ImgName)
	}

	resource := dto.EmployeeDetailResource{
		ID:        employee.ID,
		NIP:       employee.Nip,
		Position:  employee.Position,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
		Major: dto.MajorOptionResource{
			ID:   employee.Major.ID,
			Name: employee.Major.Name,
		},
		User: dto.UserDetailResource{
			ID:          employee.User.ID,
			Name:        employee.User.Name,
			Email:       employee.User.Email,
			Status:      employee.User.Status,
			Gender:      employee.User.Gender,
			Religion:    employee.User.Religion,
			BirthDate:   employee.User.BirthDate,
			BirthPlace:  employee.User.BirthPlace,
			PhoneNumber: employee.User.PhoneNumber,
			Address:     employee.User.Address,
			Nationality: employee.User.Nationality,
			Avatar:      avatarURL,
		},
	}

	helper.SuccessResponse(c, http.StatusOK, "Employee found", resource)
}

func (h *EmployeeHandler) FindAllAsOptions(c *gin.Context) {
	position := c.Query("position")
	majorId := c.Query("major_id")
	studyProgramId := c.Query("study_program_id")

	employees, err := h.useCase.FindAllAsOptions(position, majorId, studyProgramId)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch employee options", err)
		return
	}

	var optionResource []dto.Option
	for _, emp := range *employees {
		optionResource = append(optionResource, dto.Option{
			Label: emp.Name,
			Value: emp.ID,
		})
	}

	helper.SuccessResponse(c, http.StatusOK, "Employee options fetched successfully", optionResource)
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var payload dto.StoreEmployeeDTO
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

	employee, err := h.useCase.Create(&payload)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			helper.ErrorResponse(c, http.StatusConflict, "Employee with this email or NIP already exists", err)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to create employee", err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "Employee created successfully", employee)
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var payload dto.UpdateEmployeeDTO
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

	employee, err := h.useCase.Update(id, &payload)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Employee not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to update employee", err)
		}
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Employee updated successfully", employee)
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.useCase.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			helper.ErrorResponse(c, http.StatusNotFound, "Employee not found", err)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete employee", err)
		}
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "Employee deleted successfully", nil)
}
