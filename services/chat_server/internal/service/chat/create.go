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
	users, err := s.userCli.GetList(ctx, usernames)
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, fmt.Errorf("cannot create empty chat")
	}

	userIDs := make([]int64, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	chat := &model.Chat{
		Name: name,
	}

	var newChatID int64
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		newChatID, err = s.chatRepo.Create(ctx, chat)
		if err != nil {
			return fmt.Errorf("cannot create chat: %v", err)
		}

		err = s.chatMemberRepo.CreateBatch(ctx, newChatID, userIDs)
		if err != nil {
			return fmt.Errorf("cannot create chat members: %v", err)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return newChatID, nil
}
