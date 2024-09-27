package chat

import (
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/client/service"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/client/service/chat/converter"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/model"
)

var _ service.ChatClient = (*chatGRPCClient)(nil)

type chatGRPCClient struct {
	client service.ChatGRPCClient
}

// NewChatClient initialises grpc client to chat service
func NewChatClient(client service.ChatGRPCClient) *chatGRPCClient {
	return &chatGRPCClient{
		client: client,
	}
}

// Create creates new chat
func (c *chatGRPCClient) Create(ctx context.Context, name string, usernames []string) (int64, error) {
	req := &chatApi.CreateRequest{
		Name:      name,
		Usernames: usernames,
	}

	res, err := c.client.Create(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %v", err)
	}

	return res.Id, nil
}

func (c *chatGRPCClient) SendMessage(ctx context.Context, message *model.Message) error {
	messageProto := converter.ToProtoFromMessage(message)
	_, err := c.client.SendMessage(ctx, &chatApi.SendMessageRequest{
		Message: messageProto,
	})

	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	return nil
}

func (c *chatGRPCClient) Connect(ctx context.Context, chatID int64, username string) (chan *model.Message, error) {
	stream, err := c.client.Connect(ctx, &chatApi.ConnectChatRequest{ChatId: chatID, Username: username})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to chat: %v", err)
	}

	messageChannel := make(chan *model.Message, 100)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(messageChannel)
				return
			}

			if err != nil {
				close(messageChannel)
				logger.Error("error receiving message from chat", zap.Error(err))
				return
			}

			message := converter.FromProtoToMessage(resp)
			messageChannel <- message
		}
	}()

	return messageChannel, nil
}
