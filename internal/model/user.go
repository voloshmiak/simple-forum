package model

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	Role         string    `json:"role"`
}

type AuthorizedUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
