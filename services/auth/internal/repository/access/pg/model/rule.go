package model

// RoleAccessRule represent role access to endpoint
type RoleAccessRule struct {
	Role     int    `json:"role"`
	Endpoint string `json:"endpoint"`
}
