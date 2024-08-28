package access

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"

	accessApi "github.com/Genvekt/cli-chat/libraries/api/access/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
)

const (
	authorisationHeader = "authorization"
	authPrefix          = "Bearer "
)

var _ service.AccessClient = (*accessClient)(nil)

type accessClient struct {
	client service.AccessGrpcClient
}

// NewAccessClient initialises grpc client to access server
func NewAccessClient(client service.AccessGrpcClient) *accessClient {
	return &accessClient{
		client: client,
	}
}

func (c *accessClient) Check(ctx context.Context, endpoint string) (bool, error) {
	accessToken, err := c.getAccessTokenFromCtx(ctx)
	if err != nil {
		return false, err
	}

	ctx = c.putAccessTokenToCtx(ctx, accessToken)

	_, err = c.client.Check(ctx, &accessApi.CheckRequest{
		EndpointAddress: endpoint,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *accessClient) getAccessTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata not provided")
	}
	authHeader := md.Get(authorisationHeader)
	if len(authHeader) == 0 {
		return "", fmt.Errorf("%s header not provided", authorisationHeader)
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return "", fmt.Errorf("invalid %s header format", authorisationHeader)
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	return accessToken, nil
}

func (c *accessClient) putAccessTokenToCtx(ctx context.Context, accessToken string) context.Context {
	md := metadata.New(map[string]string{authorisationHeader: authPrefix + accessToken})
	return metadata.NewOutgoingContext(ctx, md)
}
