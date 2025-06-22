package dto

// LoginRequestDTO is used for the login request body.
type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponseDTO is the successful login response.
type LoginResponseDTO struct {
	Token string        `json:"token"`
	User  UserLoginInfo `json:"user"`
}

// UserLoginInfo contains the user data returned on login.
type UserLoginInfo struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
	Permissions []string `json:"permissions"`
}