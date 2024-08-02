package chat

import (
	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/db"
	serviceCli "github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/service"
)

var _ service.ChatService = (*chatService)(nil)

type chatService struct {
	chatRepo       repository.ChatRepository
	chatMemberRepo repository.ChatMemberRepository
	messageRepo    repository.MessageRepository
	userCli        serviceCli.AuthClient
	txManager      db.TxManager
}

// NewChatService initialises service layer for chat business logic
func NewChatService(
	chatRepo repository.ChatRepository,
	chatMemberRepo repository.ChatMemberRepository,
	messageRepository repository.MessageRepository,
	userCli serviceCli.AuthClient,
	txManager db.TxManager,
) *chatService {
	return &chatService{
		chatRepo:       chatRepo,
		chatMemberRepo: chatMemberRepo,
		messageRepo:    messageRepository,
		userCli:        userCli,
		txManager:      txManager,
	}
}
