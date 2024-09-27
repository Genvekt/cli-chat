package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
)

// SendMessage sends message to chat
func (s *chatService) SendMessage(ctx context.Context, message *model.Message) error {
	isMember, err := s.isMember(ctx, message.ChatID, message.SenderID)
	if err != nil {
		return fmt.Errorf("failed to check the membership of user in chat: %v", err)
	}

	// check user membership
	if !isMember {
		return fmt.Errorf("user is not a member of a chat")
	}

	// save message to storage
	err = s.messageRepo.Create(ctx, message)
	if err != nil {
		return fmt.Errorf("cannot create message: %v", err)
	}

	// send to active chat. Do nothing for inactive chats
	s.mxChat.Lock()
	if chatConnection, isActive := s.chatConnections[message.ChatID]; isActive {
		chatConnection.Buffer <- message
	}
	s.mxChat.Unlock()

	return nil
}

// isMember checks that user is a member of a chat
func (s *chatService) isMember(ctx context.Context, chatID int64, userID int64) (bool, error) {
	_, err := s.chatMemberRepo.Get(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrChatMemberNotFound) {
			return false, nil
		}
		return false, err

	}
	return true, nil
}
