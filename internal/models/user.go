package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	Role         string `json:"role"`
}

func NewUser() *User {
	return &User{}
}
