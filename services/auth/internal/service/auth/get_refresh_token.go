package auth

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// GetRefreshToken creates new refresh token based on valid refresh token
func (s *authService) GetRefreshToken(ctx context.Context, oldRefreshTocken string) (string, error) {
	claims, err := s.refreshTokenProvider.Verify(ctx, oldRefreshTocken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %v", err)
	}

	refreshToken, err := s.refreshTokenProvider.Generate(ctx, &model.User{
		Name: claims.Username,
		Role: claims.Role,
	})
	if err != nil {
		return "", fmt.Errorf("cannot generate refresh token: %v", err)
	}

	return refreshToken, nil
}
