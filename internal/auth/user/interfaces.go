package user

type UserRepository interface {
	Create(user *CreateUser) (*User, error)
	GetByUsername(username string) (*User, error)
	GetById(id int64) (*User, error)
	DeleteById(id int64) error
}
