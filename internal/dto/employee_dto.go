package dto

import (
	"mime/multipart"
	"time"
)

// Position values based on PositionEnum in Laravel
const (
	PositionLecturer = "LECTURER"
	PositionStaff    = "STAFF"
)

// StoreEmployeeDTO combines user and employee data for creation
// Based on StoreEmployeeRequest.php
type StoreEmployeeDTO struct {
	MajorID     *string               `form:"m_major_id" binding:"omitempty,uuid"`
	NIP         string                `form:"nip" binding:"required,max=255"`
	Position    string                `form:"position" binding:"required,oneof=DOSEN TEKNISI ADMINISTRASI"`
	Name        string                `form:"name" binding:"required,max=255"`
	Email       string                `form:"email" binding:"required,email,max=255"`
	Gender      *string               `form:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Religion    *string               `form:"religion" binding:"omitempty,oneof=ISLAM CHRISTIANITY CATHOLIC HINDUISM BUDDHISM CONFUCIANISM OTHER"`
	BirthDate   *string               `form:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	BirthPlace  *string               `form:"birth_place" binding:"omitempty,max=255"`
	Address     *string               `form:"address" binding:"omitempty,max=255"`
	PhoneNumber *string               `form:"phone_number" binding:"omitempty,max=20"`
	Nationality *string               `form:"nationality" binding:"omitempty,max=255"`
	Avatar      *multipart.FileHeader `form:"avatar" binding:"-"`
}

// UpdateEmployeeDTO for updates
type UpdateEmployeeDTO struct {
	MajorID     *string               `form:"m_major_id" binding:"omitempty,uuid"`
	NIP         string                `form:"nip" binding:"required,max=255"`
	Position    string                `form:"position" binding:"required,oneof=DOSEN TEKNISI ADMINISTRASI"`
	Name        string                `form:"name" binding:"required,max=255"`
	Email       string                `form:"email" binding:"required,email,max=255"`
	Gender      *string               `form:"gender" binding:"omitempty,oneof=MALE FEMALE"`
	Religion    *string               `form:"religion" binding:"omitempty,oneof=ISLAM CHRISTIANITY CATHOLIC HINDUISM BUDDHISM CONFUCIANISM OTHER"`
	BirthDate   *string               `form:"birth_date" binding:"omitempty,datetime=2006-01-02"`
	BirthPlace  *string               `form:"birth_place" binding:"omitempty,max=255"`
	Address     *string               `form:"address" binding:"omitempty,max=255"`
	PhoneNumber *string               `form:"phone_number" binding:"omitempty,max=20"`
	Nationality *string               `form:"nationality" binding:"omitempty,max=255"`
	Avatar      *multipart.FileHeader `form:"avatar" binding:"-"`
}

type EmployeeResource struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	NIP      string `json:"nip"`
	Position string `json:"position"`
	Avatar   string `json:"avatar,omitempty"`
}

type UserDetailResource struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Status      string     `json:"status"`
	Gender      *string    `json:"gender"`
	Religion    *string    `json:"religion"`
	BirthDate   *time.Time `json:"birth_date"`
	BirthPlace  *string    `json:"birth_place"`
	PhoneNumber *string    `json:"phone_number"`
	Address     *string    `json:"address"`
	Nationality *string    `json:"nationality"`
	Avatar      string     `json:"avatar,omitempty"`
}

type EmployeeDetailResource struct {
	ID        string              `json:"id"`
	NIP       string              `json:"nip"`
	Position  string              `json:"position"`
	User      UserDetailResource  `json:"user"`
	Major     MajorOptionResource `json:"major"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type EmployeeDetailInfoDTO struct {
	ID               string     `json:"id"`
	NIP              string     `json:"nip"`
	Position         string     `json:"position"`
	MajorID          *string    `json:"m_major_id"`
	MajorName        *string    `json:"major_name"`
	StudyProgramID   *string    `json:"m_study_program_id"`
	StudyProgramName *string    `json:"study_program_name"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}
