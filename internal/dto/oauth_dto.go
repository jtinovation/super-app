package dto

import "time"

type StoreOauthCodeDTO struct {
	Code        string
	ClientID    string
	UserSub     LoginResponseDTO
	RedirectURI string
	ExpiresAt   time.Time
}

type LoginRequestFormDTO struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
	ReturnTo string `form:"return_to"`
}
