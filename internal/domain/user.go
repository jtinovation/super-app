package domain

type AuthRepository interface {
	FindByEmail(email string) (*User, error)
}