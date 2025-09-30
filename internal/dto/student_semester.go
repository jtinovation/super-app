package dto

type StudentSemesterResource struct {
	ID        string `json:"id"`
	Year      int    `json:"year"`
	Semester  string `json:"semester"`
	Class     string `json:"class"`
	SessionId string `json:"session_id"`
	Session   string `json:"session"`
}

type StudentSemesterDTO struct {
	ID         string `json:"id"`
	StudentID  string `json:"student_id"`
	SemesterID string `json:"semester_id"`
	Class      string `json:"class"`
	IsActive   bool   `json:"is_active"`
}
