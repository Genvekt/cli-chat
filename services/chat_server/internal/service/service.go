package service

import (
	"context"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// ChatService performs business logic related to chats
type ChatService interface {
	Create(ctx context.Context, name string, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message *model.Message) error
	Connect(ctx context.Context, id int64, username string) (chan *model.Message, error)
}
