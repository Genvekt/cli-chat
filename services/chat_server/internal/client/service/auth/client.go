package auth

import (
	"context"
	"fmt"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service/auth/converter"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

var _ service.AuthClient = (*authGRPCClient)(nil)

type authGRPCClient struct {
	client service.UserGrpcClient
}

// NewAuthClient initialises grpc client to auth service
func NewAuthClient(client service.UserGrpcClient) *authGRPCClient {
	return &authGRPCClient{
		client: client,
	}
}

// GetList retrieves users by their usernames
func (c *authGRPCClient) GetList(ctx context.Context, usernames []string) ([]*model.User, error) {
	req := &userApi.GetListRequest{Names: usernames}

	res, err := c.client.GetList(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}

	return converter.ToUsersFromRepo(res.Users), nil
}
