package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/utils"
)

func (s *CliService) Connect(ctx context.Context, chatID int64) error {
	accessToken, err := s.login(ctx)
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	ctxWithToken := utils.PutAccessTokenToCtx(ctx, accessToken)
	userID, err := utils.GetUserIdFromToken(accessToken)
	if err != nil {
		return fmt.Errorf("missing current user id in token")
	}

	msgChannel, err := s.chatClient.Connect(ctxWithToken, chatID, s.profileConfig.Username())
	if err != nil {
		return fmt.Errorf("failed to connect to chat: %v", err)
	}

	fmt.Printf("Connected to chat with ID %d\n", chatID)

	for {
		select {
		case msg, ok := <-msgChannel:
			if !ok {
				return nil
			}

			if msg.SenderID == userID {
				fmt.Printf("[%s] YOU: %s\n", msg.Timestamp.Format(time.DateTime), msg.Text)
			} else {
				fmt.Printf("[%s] %d: %s\n", msg.Timestamp.Format(time.DateTime), msg.SenderID, msg.Text)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
