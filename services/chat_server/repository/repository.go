package repository

import "github.com/Genvekt/cli-chat/services/chat-server/model"

// ChatRepository manages chats in some data source
type ChatRepository interface {
	Create(chat *model.Chat) (*model.Chat, error)
	Delete(id int64) error
}
