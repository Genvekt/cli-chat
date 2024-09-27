package auth

import (
	"context"
	"fmt"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/client/service"
)

var _ service.AuthClient = (*authGRPCClient)(nil)

type authGRPCClient struct {
	client service.AuthGRPCClient
}

// NewAuthClient initialises grpc client to auth service
func NewAuthClient(client service.AuthGRPCClient) *authGRPCClient {
	return &authGRPCClient{
		client: client,
	}
}

// Login retrieves refresh token for user
func (c *authGRPCClient) Login(ctx context.Context, username string, password string) (string, error) {
	req := &authApi.LoginRequest{
		Username: username,
		Password: password,
	}

	res, err := c.client.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to login: %v", err)
	}

	return res.RefreshToken, nil
}

// GetAccessToken retrieves access token by refresh token
func (c *authGRPCClient) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	resp, err := c.client.GetAccessToken(ctx, &authApi.GetAccessTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access token: %v", err)
	}

	return resp.AccessToken, nil
}
