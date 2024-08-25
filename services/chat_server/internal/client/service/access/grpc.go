package access

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	accessApi "github.com/Genvekt/cli-chat/libraries/api/access/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
)

var _ service.AccessGrpcClient = (*accessGRPCClient)(nil)

type accessGRPCClient struct {
	client accessApi.AccessV1Client
}

// NewAccessGrpcClient initialises grpc client to access service
func NewAccessGrpcClient(conn *grpc.ClientConn) *accessGRPCClient {
	return &accessGRPCClient{
		client: accessApi.NewAccessV1Client(conn),
	}
}

func (c *accessGRPCClient) Check(ctx context.Context, req *accessApi.CheckRequest) (*emptypb.Empty, error) {
	LogRequest(ctx, req)

	resp, err := c.client.Check(ctx, req)
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
