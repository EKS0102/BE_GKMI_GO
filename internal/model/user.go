package model

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	IsActive     bool   `json:"is_active"`
	Email        string `json:"email"`
}
