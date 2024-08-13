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
		txChatID, txErr := s.chatRepo.Create(ctx, chat)
		if txErr != nil {
			return fmt.Errorf("cannot create chat: %v", txErr)
		}

		txErr = s.chatMemberRepo.CreateBatch(ctx, txChatID, userIDs)
		if txErr != nil {
			return fmt.Errorf("cannot create chat members: %v", txErr)
		}

		newChatID = txChatID

		return nil
	})
	if err != nil {
		return 0, err
	}

	return newChatID, nil
}
