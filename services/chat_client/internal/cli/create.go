package cli

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/utils"
)

func (s *CliService) CreateChat(ctx context.Context, name string, usernames []string) error {
	token, err := s.login(ctx)
	if err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	ctxWithToken := utils.PutAccessTokenToCtx(ctx, token)

	chatID, err := s.chatClient.Create(ctxWithToken, name, usernames)
	if err != nil {
		return err
	}

	fmt.Printf("created chat '%s' (id: %d)", name, chatID)
	return nil
}
