package model

import "time"

// User service model
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Role         int       `json:"role"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserUpdateDTO is collection of values for user update
type UserUpdateDTO struct {
	ID    int64
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *int    `json:"role"`
}

// UserFilters is a collection of AND filters that may be used to query users
type UserFilters struct {
	Names []string `json:"names"`
}
