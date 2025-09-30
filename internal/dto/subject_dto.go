// internal/dto/subject_dto.go

package dto

type StoreSubjectDTO struct {
	StudyProgramID string  `json:"study_program_id" binding:"required,uuid"`
	Code           string  `json:"code" binding:"required,max=255"`
	Name           string  `json:"name" binding:"required,max=255"`
	Status         *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type UpdateSubjectDTO struct {
	StudyProgramID string  `json:"study_program_id" binding:"required,uuid"`
	Code           string  `json:"code" binding:"required,max=255"`
	Name           string  `json:"name" binding:"required,max=255"`
	Status         *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}

type SubjectResource struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	StudyProgramName string `json:"study_program_name"`
	StudyProgramID   string `json:"study_program_id"`
}

type LectureMappingDTO struct {
	SubjectSemesterID string   `json:"subject_semester_id" binding:"required,uuid"`
	LectureIDs        []string `json:"lecture_ids" binding:"omitempty,dive,uuid"`
}

type SettingLectureOnSubjectDTO struct {
	Data []LectureMappingDTO `json:"data" binding:"required,min=1,dive"`
}

type LectureUserResource struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type LectureResource struct {
	ID               string              `json:"id"`
	MajorID          *string             `json:"major_id"`
	User             LectureUserResource `json:"user"`
	SubjectLectureID string              `json:"subject_lecture_id"`
}

type LectureOnSubjectSubjectResource struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	StudyProgramID string `json:"study_program_id"`
}

type LectureOnSubjectResource struct {
	ID         string                          `json:"id"`
	SemesterID string                          `json:"semester_id"`
	Subject    LectureOnSubjectSubjectResource `json:"subject"`
	Lectures   []LectureResource               `json:"lectures"`
}
