package dto

import (
	"mime/multipart"
	"time"
)

type StoreStudentDTO struct {
	Name           string  `form:"name" binding:"required,max=255"`
	Email          string  `form:"email" binding:"required,email,max=255"`
	NIM            string  `form:"nim" binding:"required,max=255"`
	Generation     *int    `form:"generation" binding:"omitempty,number"`
	Gender         *string `form:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Religion       *string `form:"religion" binding:"omitempty,oneof=ISLAM CHRISTIANITY CATHOLIC HINDUISM BUDDHISM CONFUCIANISM OTHER"`
	BirthPlace     *string `form:"birth_place" binding:"omitempty,max=255"`
	BirthDate      *string `form:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	TuitionFee     *int    `form:"tuition_fee" binding:"omitempty,number"`
	TuitionMethod  *string `form:"tuition_method" binding:"omitempty,max=255"`
	Address        *string `form:"address" binding:"omitempty,max=255"`
	PhoneNumber    *string `form:"phone_number" binding:"omitempty,max=20"`
	Nationality    *string `form:"nationality" binding:"omitempty,max=255"`
	SemesterId     string  `form:"semester_id" binding:"required,uuid"`
	StudyProgramID string  `form:"study_program_id" binding:"required,uuid"`
	Class          string  `form:"class" binding:"required,max=255"`

	Avatar *multipart.FileHeader `form:"avatar" binding:"-"`
}

type UpdateStudentDTO struct {
	Name           string  `form:"name" binding:"required,max=255"`
	Email          string  `form:"email" binding:"required,email,max=255"`
	NIM            string  `form:"nim" binding:"required,max=255"`
	Generation     *int    `form:"generation" binding:"omitempty,number"`
	Gender         *string `form:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Religion       *string `form:"religion" binding:"omitempty,oneof=ISLAM CHRISTIANITY CATHOLIC HINDUISM BUDDHISM CONFUCIANISM OTHER"`
	BirthPlace     *string `form:"birth_place" binding:"omitempty,max=255"`
	BirthDate      *string `form:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	TuitionFee     *int    `form:"tuition_fee" binding:"omitempty,number"`
	TuitionMethod  *string `form:"tuition_method" binding:"omitempty,max=255"`
	Address        *string `form:"address" binding:"omitempty,max=255"`
	PhoneNumber    *string `form:"phone_number" binding:"omitempty,max=20"`
	Nationality    *string `form:"nationality" binding:"omitempty,max=255"`
	Status         *string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	StudyProgramID string  `form:"study_program_id" binding:"required,uuid"`

	Avatar *multipart.FileHeader `form:"avatar" binding:"-"`
}

type StudentResource struct {
	ID           string                     `json:"id"`
	UserID       string                     `json:"user_id"`
	NIM          string                     `json:"nim"`
	Name         string                     `json:"name"`
	Generation   *int                       `json:"generation"`
	Class        string                     `json:"class,omitempty"`
	StudyProgram StudyProgramOptionResource `json:"study_program"`
	Major        MajorOptionResource        `json:"major"`
	Avatar       string                     `json:"avatar,omitempty"`
}

type StudentDetailResource struct {
	ID             string                    `json:"id"`
	NIM            string                    `json:"nim"`
	Generation     *int                      `json:"generation"`
	TuitionFee     *int                      `json:"tuition_fee"`
	TuitionMethod  *string                   `json:"tuition_method"`
	StudyProgramId string                    `json:"study_program_id"`
	MajorId        string                    `json:"major_id"`
	User           UserResource              `json:"user"`
	Semesters      []StudentSemesterResource `json:"semesters"`
}

type StudentDetailInfoDTO struct {
	ID               string                `json:"id"`
	NIM              string                `json:"nim"`
	Generation       *int                  `json:"generation"`
	MajorID          *string               `json:"m_major_id"`
	MajorName        *string               `json:"major_name"`
	StudyProgramID   *string               `json:"m_study_program_id"`
	StudyProgramName *string               `json:"study_program_name"`
	StudentSemesters *[]StudentSemesterDTO `json:"student_semester"`
	CreatedAt        *time.Time            `json:"created_at,omitempty"`
	UpdatedAt        *time.Time            `json:"updated_at,omitempty"`
}
