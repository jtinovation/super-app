package dto

type StoreOauthClientDTO struct {
	Name     string `json:"name" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
	Redirect string `json:"redirect" binding:"required,url"`
}

type UpdateOauthClientDTO struct {
	Name     string `json:"name" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
	Redirect string `json:"redirect" binding:"required,url"`
}

type OauthClientResource struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Redirect string `json:"redirect"`
}
