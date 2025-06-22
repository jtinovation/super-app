package dto

type StoreStudyProgramDTO struct {
	MajorID string `json:"major_id" binding:"required,uuid"`
	Code    string `json:"code" binding:"required,max=255"`
	Name    string `json:"name" binding:"required,max=255"`
}

type UpdateStudyProgramDTO struct {
	MajorID string `json:"major_id" binding:"required,uuid"`
	Code    string `json:"code" binding:"required,max=255"`
	Name    string `json:"name" binding:"required,max=255"`
}

type StudyProgramResource struct {
	ID    string              `json:"id"`
	Code  string              `json:"code"`
	Name  string              `json:"name"`
	Major MajorOptionResource `json:"major"`
}

type StudyProgramOptionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
