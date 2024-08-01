package chat

import (
	"context"
	"errors"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
)

// IsMember checks that user is a member of a chat
func (s *chatService) IsMember(ctx context.Context, chatID int64, userID int64) (bool, error) {
	_, err := s.chatRepo.GetChatMember(ctx, chatID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrChatMemberNotFound) {
			return false, nil
		}
		return false, err

	}
	return true, nil
}
