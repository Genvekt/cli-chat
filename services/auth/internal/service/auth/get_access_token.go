package auth

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// GetAccessToken creates new refresh token based on valid refresh token
func (s *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.refreshTokenProvider.Verify(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %v", err)
	}

	accessToken, err := s.accessTokenProvider.Generate(ctx, &model.User{
		ID:   claims.ID,
		Name: claims.Username,
		Role: claims.Role,
	})
	if err != nil {
		return "", fmt.Errorf("cannot generate access token: %v", err)
	}

	return accessToken, nil
}
