package auth

import "time"

type User struct {
	Id           int64
	CreatedAt    time.Time
	FirstName    string
	LastName     string
	Username     string
	Role         string
	PasswordHash []byte
	PasswordSalt []byte
}

type UserOut struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
}

type CreateUser struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Role      string `json:"role" binding:"required"`
	Password  string `json:"password" binding:"password"`
}
