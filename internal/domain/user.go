package domain

type AuthRepository interface {
	FindByEmail(email string) (*User, error)
}

type UserRepository interface {
	FindByID(id string) (*User, error)
	Create(user *User) (*User, error)
	Update(id string, user *User) (*User, error)
}
