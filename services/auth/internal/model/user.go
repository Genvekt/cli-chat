package model

import "time"

// User service model
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserUpdateDTO is collection of values for user update
type UserUpdateDTO struct {
	ID    int64
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *int    `json:"role"`
}
