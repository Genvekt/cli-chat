package auth

import (
	"context"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
)

// GetRefreshToken provides refresh token for valid refresh token
func (s *Service) GetRefreshToken(ctx context.Context, req *authApi.GetRefreshTokenRequest) (*authApi.GetRefreshTokenResponse, error) {
	refreshToken, err := s.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authApi.GetRefreshTokenResponse{
		RefreshToken: refreshToken,
	}, nil
}
