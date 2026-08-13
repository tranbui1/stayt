package models

type User struct {
	ID           int    // ID of user
	UserName     string // Username
	CreatedAt    string // Date account created
	Email        string // Email used to create account
	PasswordHash string
}
