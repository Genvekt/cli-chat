package cli

import (
	clientService "github.com/Genvekt/cli-chat/services/chat-client/internal/client/service"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/config"
)

type CliService struct {
	profileConfig config.ProfileConfig
	chatClient    clientService.ChatClient
	authClient    clientService.AuthClient
}

func NewChatCliService(
	profileConfig config.ProfileConfig,
	chatClient clientService.ChatClient,
	authClient clientService.AuthClient,
) *CliService {
	return &CliService{
		profileConfig: profileConfig,
		chatClient:    chatClient,
		authClient:    authClient,
	}
}
