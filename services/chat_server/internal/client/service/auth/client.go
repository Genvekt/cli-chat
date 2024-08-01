package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service/auth/converter"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

var _ service.AuthClient = (*authGRPCClient)(nil)

type authGRPCClient struct {
	// Not sure that there may be several instances of auth services,
	// wrapper is done to understand the concept better
	client service.UserGrpcClient
	conn   *grpc.ClientConn
}

// NewAuthClient initialises grpc client to auth service
func NewAuthClient(address string) (*authGRPCClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	return &authGRPCClient{
		client: NewAuthGrpcClient(conn),
		conn:   conn,
	}, nil
}

// Query retrieves users by their usernames
func (c *authGRPCClient) Query(ctx context.Context, usernames []string) ([]*model.User, error) {
	req := &userApi.QueryRequest{Names: usernames}

	res, err := c.client.Query(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}

	var users []*model.User
	for _, user := range res.Users {
		users = append(users, converter.ToUserFromRepo(user))
	}

	return users, nil
}

// Close closes grpc connection
func (c *authGRPCClient) Close() error {
	return c.conn.Close()
}
