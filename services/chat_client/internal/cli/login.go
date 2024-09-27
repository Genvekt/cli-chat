package cli

import (
	"context"
)

func (s *CliService) login(ctx context.Context) (string, error) {
	refreshToken, err := s.authClient.Login(ctx, s.profileConfig.Username(), s.profileConfig.Password())
	if err != nil {
		return "", err
	}

	accessToken, err := s.authClient.GetAccessToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
