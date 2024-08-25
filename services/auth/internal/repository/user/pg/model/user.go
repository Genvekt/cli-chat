package user

import "time"

// User repository model
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         int       `json:"role"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Filters is a collection of AND filters that may be used to query users
type Filters struct {
	Names []string `json:"names"`
}
