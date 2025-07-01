package dto

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
