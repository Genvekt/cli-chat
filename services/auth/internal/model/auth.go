package model

import "github.com/golang-jwt/jwt"

// UserClaims is data for token
type UserClaims struct {
  jwt.StandardClaims
  ID       int64  `json:"id"`
  Username string `json:"username"`
  Role     int    `json:"role"`
}
