package dto

type StoreClassDTO struct {
	StudyProgramID string `json:"study_program_id" binding:"required,uuid"`
	Code           string `json:"code" binding:"required,max=255"`
	Name           string `json:"name" binding:"required,max=255"`
}

type UpdateClassDTO struct {
	StudyProgramID string `json:"study_program_id" binding:"required,uuid"`
	Code           string `json:"code" binding:"required,max=255"`
	Name           string `json:"name" binding:"required,max=255"`
}

type ClassResource struct {
	ID           string                     `json:"id"`
	Code         string                     `json:"code"`
	Name         string                     `json:"name"`
	StudyProgram StudyProgramOptionResource `json:"study_program"`
	Major        MajorOptionResource        `json:"major"`
}

type ClassOptionResource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
