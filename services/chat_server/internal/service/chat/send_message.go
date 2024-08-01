package chat

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// SendMessage sends message to chat
func (s *chatService) SendMessage(ctx context.Context, message *model.Message) error {
	isMember, err := s.IsMember(ctx, message.ChatID, message.SenderID)
	if err != nil {
		return fmt.Errorf("failed to check the membership of user in chat: %v", err)
	}
	if !isMember {
		return fmt.Errorf("user is not a member of a chat")
	}
	err = s.messageRepo.Create(ctx, message)
	if err != nil {
		return err
	}

	return nil
}
