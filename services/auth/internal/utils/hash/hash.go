package hash

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type hasher struct{}

// NewHasher initialises hashing utils
func NewHasher() *hasher {
	return &hasher{}
}

// HashPassword converts password string to hash
func (h *hasher) HashPassword(_ context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares password with target hash
func (h *hasher) CheckPasswordHash(_ context.Context, password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
