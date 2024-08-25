package auth

import (
	"context"
	"fmt"
)

// Login creates new refresh token based on username-password pair
func (s *authService) Login(ctx context.Context, username, password string) (string, error) {

	// Get user with provided username
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("cannot get user with username %s: %v", username, err)
	}

	// Validate password
	passwordCorrect := s.hasher.CheckPasswordHash(ctx, password, user.PasswordHash)
	if !passwordCorrect {
		return "", fmt.Errorf("invalid password")
	}

	refreshToken, err := s.refreshTokenProvider.Generate(ctx, user)
	if err != nil {
		return "", fmt.Errorf("cannot generate refresh token: %v", err)
	}

	return refreshToken, nil
}
