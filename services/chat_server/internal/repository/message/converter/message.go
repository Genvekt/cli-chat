package converter

import (
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/message/model"
)

// ToRepoFromMessage converts message service model to message repository model
func ToRepoFromMessage(message *model.Message) *repoModel.Message {
	return &repoModel.Message{
		ID:        message.ID,
		SenderID:  message.SenderID,
		ChatID:    message.ChatID,
		Content:   message.Content,
		Timestamp: message.Timestamp,
	}
}
