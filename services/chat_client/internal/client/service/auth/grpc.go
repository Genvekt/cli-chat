package auth

import (
	"context"

	"google.golang.org/grpc"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/client/service"
)

var _ service.AuthGRPCClient = (*authGrpcClientWrapper)(nil)

type authGrpcClientWrapper struct {
	client authApi.AuthV1Client
}

// NewAuthGrpcClientWrapper initialises wrapper around grpc client
func NewAuthGrpcClientWrapper(conn *grpc.ClientConn) *authGrpcClientWrapper {
	return &authGrpcClientWrapper{
		client: authApi.NewAuthV1Client(conn),
	}
}

// Login performs login request
func (c *authGrpcClientWrapper) Login(
	ctx context.Context,
	req *authApi.LoginRequest,
) (*authApi.LoginResponse, error) {

	resp, err := c.client.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *authGrpcClientWrapper) GetAccessToken(
	ctx context.Context,
	req *authApi.GetAccessTokenRequest,
) (*authApi.GetAccessTokenResponse, error) {
	resp, err := c.client.GetAccessToken(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
