package auth

import (
	"context"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
)

// GetAccessToken provides access token for valid refresh token
func (s *Service) GetAccessToken(ctx context.Context, req *authApi.GetAccessTokenRequest) (*authApi.GetAccessTokenResponse, error) {
	refreshToken, err := s.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authApi.GetAccessTokenResponse{
		AccessToken: refreshToken,
	}, nil
}
