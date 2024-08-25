package auth

import (
	"context"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
)

// Login provides refresh token for vald username-password pair
func (s *Service) Login(ctx context.Context, req *authApi.LoginRequest) (*authApi.LoginResponse, error) {
	refreshToken, err := s.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &authApi.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
