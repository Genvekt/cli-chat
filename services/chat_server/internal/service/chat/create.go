package chat

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// Create creates chat
func (s *chatService) Create(ctx context.Context, name string, usernames []string) (int64, error) {
	// Get users by their ids from other microservice
	// Done to implement service to service communication
	users, err := s.userCli.Query(ctx, usernames)
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, fmt.Errorf("cannot create empty chat")
	}

	chatMembers := make([]*model.ChatMember, 0, len(users))
	for _, user := range users {
		chatMembers = append(chatMembers, &model.ChatMember{ID: user.ID})
	}

	chat := &model.Chat{
		Name:    name,
		Members: chatMembers,
	}

	var newChatID int64
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		newChatID, err = s.chatRepo.Create(ctx, chat)
		if err != nil {

			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return newChatID, nil
}
