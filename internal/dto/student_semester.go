package dto

type StudentSemesterResource struct {
	ID        string `json:"id"`
	Year      int    `json:"year"`
	Semester  string `json:"semester"`
	Class     string `json:"class"`
	SessionId string `json:"session_id"`
	Session   string `json:"session"`
}
