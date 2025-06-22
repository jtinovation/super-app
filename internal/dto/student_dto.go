package dto

import "mime/multipart"

type StoreStudentDTO struct {
	ClassID       string  `form:"class_id" binding:"required,uuid"`
	Name          string  `form:"name" binding:"required,max=255"`
	NIM           string  `form:"nim" binding:"required,max=255"`
	Generation    *int    `form:"generation" binding:"omitempty,number"`
	Gender        *string `form:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Religion      *string `form:"religion" binding:"omitempty,oneof=ISLAM CHRISTIANITY CATHOLIC HINDUISM BUDDHISM CONFUCIANISM OTHER"`
	BirthPlace    *string `form:"birth_place" binding:"omitempty,max=255"`
	BirthDate     *string `form:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	TuitionFee    *int    `form:"tuition_fee" binding:"omitempty,number"`
	TuitionMethod *string `form:"tuition_method" binding:"omitempty,max=255"`

	Avatar *multipart.FileHeader `form:"avatar" binding:"-"`
}

type StudentResource struct {
	ID           string                     `json:"id"`
	NIM          string                     `json:"nim"`
	Name         string                     `json:"name"`
	Generation   *int                       `json:"generation"`
	Class        ClassOptionResource        `json:"class"`
	StudyProgram StudyProgramOptionResource `json:"study_program"`
	Major        MajorOptionResource        `json:"major"`
	Avatar       string                     `json:"avatar,omitempty"`
}
