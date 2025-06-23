package dto

type StoreSemesterDTO struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
	Year      int    `json:"year" binding:"required,number,gte=2000,lte=2100"`
	Semester  string `json:"semester" binding:"required,max=2"`
}

type UpdateSemesterDTO struct {
	SessionID string `json:"session_id" binding:"required,uuid"`
	Year      int    `json:"year" binding:"required,number,gte=2000,lte=2100"`
	Semester  string `json:"semester" binding:"required,max=2"`
}

type SettingSubjectSemesterDTO struct {
	SubjectIDs []string `json:"subject_ids" binding:"omitempty,dive,uuid"`
}

type SemesterResource struct {
	ID       string          `json:"id"`
	Year     int             `json:"year"`
	Semester string          `json:"semester"`
	Session  SessionResource `json:"session"`
}

type SemesterOptionResource struct {
	ID       string `json:"id"`
	Year     int    `json:"year"`
	Semester string `json:"semester"`
}
