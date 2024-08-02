package auth

import (
	"context"
	"fmt"
	"log"

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
	LogRequest(ctx, req)

	resp, err := c.client.GetList(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// LogRequest logs request sent to grpc client
func LogRequest(ctx context.Context, req interface{}) {
	log.Println(
		ctx,
		fmt.Sprintf("qrpc: %T", req),
		fmt.Sprintf("request: %+v", req),
	)
}
