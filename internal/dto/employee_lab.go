package dto

type StoreEmployeeLabDTO struct {
	LabID      string  `json:"lab_id" binding:"required,uuid"`
	EmployeeID string  `json:"employee_id" binding:"required,uuid"`
	IsHeadLab  bool    `json:"is_head_lab" binding:"required"`
	Period     *string `json:"period"`
}

type StoreEmployeeLabFromLabDTO struct {
	EmployeeID string  `json:"employee_id" binding:"required,uuid"`
	IsHeadLab  bool    `json:"is_head_lab" binding:"required"`
	Period     *string `json:"period"`
}

type UpdateEmployeeLabDTO struct {
	LabID      string  `json:"lab_id" binding:"required,uuid"`
	EmployeeID string  `json:"employee_id" binding:"required,uuid"`
	IsHeadLab  bool    `json:"is_head_lab" binding:"required"`
	Period     *string `json:"period"`
	Status     string  `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type EmployeeLabResource struct {
	ID        string            `json:"id"`
	Lab       LabOptionResource `json:"lab"`
	Employee  EmployeeResource  `json:"employee"`
	IsHeadLab bool              `json:"is_head_lab"`
	Period    *string           `json:"period"`
	Status    string            `json:"status"`
}

type EmployeeLabResourceMany struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	IsHeadLab bool    `json:"is_head_lab"`
	Period    *string `json:"period"`
	Status    string  `json:"status"`
}
