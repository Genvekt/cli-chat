package chat

import (
	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/service"
)

// Service implements chat api
type Service struct {
	chatApi.UnimplementedChatV1Server
	chatService service.ChatService
}

// NewService initialises chat api implementation
func NewService(chatService service.ChatService) *Service {
	return &Service{chatService: chatService}
}
