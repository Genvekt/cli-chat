package converter

import (
	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/model"
)

// ToChatFromProto converts user model of repository layer to user model of service layer
func ToChatFromProto(chat *chatApi.Chat) *model.Chat {
	return &model.Chat{
		ID:        chat.Id,
		Name:      chat.Info.Name,
		CreatedAt: chat.CreatedAt.AsTime(),
	}
}
