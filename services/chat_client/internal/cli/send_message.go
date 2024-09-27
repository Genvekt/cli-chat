package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/model"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/utils"
)

func (s *CliService) SendMessage(ctx context.Context, chatID int64, message string) error {
	accessToken, err := s.login(ctx)
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	ctxWithToken := utils.PutAccessTokenToCtx(ctx, accessToken)

	userID, err := utils.GetUserIdFromToken(accessToken)
	if err != nil {
		return fmt.Errorf("missing id if current user in token")
	}

	err = s.chatClient.SendMessage(ctxWithToken, &model.Message{
		ChatID:    chatID,
		SenderID:  userID,
		Text:      message,
		Timestamp: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
