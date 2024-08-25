package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	serviceClient "github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
)

// Authorization provides authorization interceptor
type Authorization struct {
	accessClient serviceClient.AccessClient
}

// NewAuthorization initialises authorisation instance
func NewAuthorization(accessClient serviceClient.AccessClient) *Authorization {
	return &Authorization{
		accessClient: accessClient,
	}
}

// Interceptor checks client access rights before handle gRPC request
func (a *Authorization) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ok, err := a.accessClient.Check(ctx, info.FullMethod)
	if err != nil || !ok {
		return nil, fmt.Errorf("access denied: %v", err)
	}

	return handler(ctx, req)
}
