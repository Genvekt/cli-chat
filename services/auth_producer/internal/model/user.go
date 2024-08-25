package model

// User service model
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
	Password string `json:"password"`
}
