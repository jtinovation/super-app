package dto

type StoreSessionDTO struct {
	Session string `json:"session" binding:"required,max=255"`
}

type UpdateSessionDTO struct {
	Session string `json:"session" binding:"required,max=255"`
}

type SessionResource struct {
	ID      string `json:"id"`
	Session string `json:"session"`
}
