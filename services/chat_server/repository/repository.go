package repository

import "github.com/Genvekt/cli-chat/services/chat-server/model"

type ChatRepository interface {
	Create(chat *model.Chat) (*model.Chat, error)
	Delete(id int64) error
}
