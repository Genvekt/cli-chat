package repository

import (
	"context"
	"errors"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// ErrChatMemberNotFound indicates chat member not found in data source
var ErrChatMemberNotFound = errors.New("chat member not found")

// ChatRepository manages chats in some data source
type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// ChatMemberRepository manages chat members in some data source
type ChatMemberRepository interface {
	CreateBatch(ctx context.Context, chatID int64, userIDs []int64) error
	Get(ctx context.Context, chatID int64, userID int64) (*model.ChatMember, error)
}

// MessageRepository manages messages in some data source
type MessageRepository interface {
	Create(ctx context.Context, message *model.Message) error
}
