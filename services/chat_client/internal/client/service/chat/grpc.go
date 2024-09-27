package chat

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/client/service"
)

var _ service.ChatGRPCClient = (*chatGrpcClientWrapper)(nil)

type chatGrpcClientWrapper struct {
	client chatApi.ChatV1Client
}

// NewChatGrpcClientWrapper initialises wrapper around grpc client
func NewChatGrpcClientWrapper(conn *grpc.ClientConn) *chatGrpcClientWrapper {
	return &chatGrpcClientWrapper{
		client: chatApi.NewChatV1Client(conn),
	}
}

// Create performs Create chat request
func (c *chatGrpcClientWrapper) Create(
	ctx context.Context,
	req *chatApi.CreateRequest,
) (*chatApi.CreateResponse, error) {

	resp, err := c.client.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *chatGrpcClientWrapper) SendMessage(
	ctx context.Context,
	req *chatApi.SendMessageRequest,
) (*emptypb.Empty, error) {

	resp, err := c.client.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *chatGrpcClientWrapper) Connect(
	ctx context.Context,
	req *chatApi.ConnectChatRequest,
) (chatApi.ChatV1_ConnectChatClient, error) {
	resp, err := c.client.ConnectChat(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
