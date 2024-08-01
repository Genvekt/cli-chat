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

// Query performs QueryRequest
func (c *authGrpcClient) Query(
	ctx context.Context,
	req *userApi.QueryRequest,
) (*userApi.QueryResponse, error) {
	LogRequest(ctx, req)
	return c.client.Query(ctx, req)
}

// LogRequest logs request sent to grpc client
func LogRequest(ctx context.Context, req interface{}) {
	log.Println(
		ctx,
		fmt.Sprintf("qrpc: %T", req),
		fmt.Sprintf("request: %+v", req),
	)
}
