package dto

import "time"

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	Token string        `json:"token"`
	User  UserLoginInfo `json:"user"`
}

type UserLoginInfo struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Email            string   `json:"email"`
	IsChangePassword bool     `json:"is_change_password"`
	Roles            []string `json:"roles"`
	Permissions      []string `json:"permissions"`
}

type ForgotPasswordRequestDTO struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequestDTO struct {
	Token           string `json:"token" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

type EmailTemplateAuthDataDto struct {
	Name    string
	Link    string
	LogoURL string
}

type UserDetailInfoDTO struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Email            string                 `json:"email"`
	EmailVerifiedAt  *time.Time             `json:"email_verified_at,omitempty"`
	Status           string                 `json:"status"`
	Gender           *string                `json:"gender,omitempty"`
	Religion         *string                `json:"religion,omitempty"`
	BirthPlace       *string                `json:"birth_place,omitempty"`
	BirthDate        *time.Time             `json:"birth_date,omitempty"`
	PhoneNumber      *string                `json:"phone_number,omitempty"`
	Address          *string                `json:"address,omitempty"`
	Nationality      *string                `json:"nationality,omitempty"`
	ImgPath          *string                `json:"img_path,omitempty"`
	ImgName          *string                `json:"img_name,omitempty"`
	IsChangePassword bool                   `json:"is_change_password"`
	Roles            []string               `json:"roles"`
	Permissions      []string               `json:"permissions"`
	CreatedAt        *time.Time             `json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty"`
	EmployeeDetail   *EmployeeDetailInfoDTO `json:"employee_detail,omitempty"`
	StudentDetail    *StudentDetailInfoDTO  `json:"student_detail,omitempty"`
}
