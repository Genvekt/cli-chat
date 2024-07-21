package repository

import (
	"context"

	"github.com/Genvekt/cli-chat/services/chat-server/model"
)

// ChatRepository manages chats in some data source
type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) error
	Delete(ctx context.Context, id int64) error
}
