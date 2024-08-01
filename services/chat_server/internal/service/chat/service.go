package chat

import (
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/db"
	serviceCli "github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/service"
)

var _ service.ChatService = (*chatService)(nil)

type chatService struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
	userCli     serviceCli.AuthClient
	txManager   db.TxManager
}

// NewChatService initialises service layer for chat business logic
func NewChatService(
	chatRepo repository.ChatRepository,
	messageRepository repository.MessageRepository,
	userCli serviceCli.AuthClient,
	txManager db.TxManager,
) *chatService {
	return &chatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepository,
		userCli:     userCli,
		txManager:   txManager,
	}
}
