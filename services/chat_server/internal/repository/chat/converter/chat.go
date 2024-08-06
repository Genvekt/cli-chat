package converter

import (
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat/model"
)

// ToRepoFromChat converts service model chat to repository model chat
func ToRepoFromChat(chat *model.Chat) *repoModel.Chat {
	chatRepo := &repoModel.Chat{
		ID:        chat.ID,
		Name:      chat.Name,
		CreatedAt: chat.CreatedAt,
	}

	return chatRepo
}
