package dto

type StoreMajorDTO struct {
	Code string `json:"code" binding:"required,max=255"`
	Name string `json:"name" binding:"required,max=255"`
}

type UpdateMajorDTO struct {
	Code string `json:"code" binding:"required,max=255"`
	Name string `json:"name" binding:"required,max=255"`
}

type MajorResource struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type MajorOptionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
