package domain

type AuthRepository interface {
	FindByEmail(email string) (*User, error)
}

type UserRepository interface {
	Create(user *User) (*User, error)
	Update(id string, user *User) (*User, error)
}
