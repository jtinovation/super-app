package dto

type StoreLabDTO struct {
	MajorID     string                       `json:"major_id" binding:"required,uuid"`
	Code        string                       `json:"code" binding:"required,max=255"`
	Name        string                       `json:"name" binding:"required,max=255"`
	EmployeeLab []StoreEmployeeLabFromLabDTO `json:"employees"`
}

type UpdateLabDTO struct {
	MajorID string `json:"major_id" binding:"required,uuid"`
	Code    string `json:"code" binding:"required,max=255"`
	Name    string `json:"name" binding:"required,max=255"`
}

type LabResource struct {
	ID          string                    `json:"id"`
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Major       MajorOptionResource       `json:"major"`
	EmployeeLab []EmployeeLabResourceMany `json:"employees"`
}

type LabOptionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
