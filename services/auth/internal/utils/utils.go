package utils

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Hasher provides hashing utils
type Hasher interface {
	HashPassword(ctx context.Context, password string) (string, error)
	CheckPasswordHash(ctx context.Context, password, hash string) bool
}

// TokenProvider provides token utils
type TokenProvider interface {
	Generate(ctx context.Context, user *model.User) (string, error)
	Verify(ctx context.Context, token string) (*model.UserClaims, error)
}
