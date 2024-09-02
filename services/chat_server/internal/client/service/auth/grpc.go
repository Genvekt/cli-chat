package auth

import (
	"context"

	"google.golang.org/grpc"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
)

var _ service.UserGrpcClient = (*authGrpcClient)(nil)

type authGrpcClient struct {
	client userApi.UserV1Client
}

// NewAuthGrpcClient initialises wrapper around grpc client
func NewAuthGrpcClient(conn *grpc.ClientConn) *authGrpcClient {
	return &authGrpcClient{
		client: userApi.NewUserV1Client(conn),
	}
}

// GetList performs GetListRequest
func (c *authGrpcClient) GetList(
	ctx context.Context,
	req *userApi.GetListRequest,
) (*userApi.GetListResponse, error) {

	resp, err := c.client.GetList(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
