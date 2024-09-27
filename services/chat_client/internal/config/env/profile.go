package env

import (
	"fmt"
	"os"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/config"
)

const (
	usernameEnv = "CLI_CHAT_USERNAME"
	passwordEnv = "CLI_CHAT_PASSWORD"
)

var _ config.ProfileConfig = (*profileConfigEnv)(nil)

type profileConfigEnv struct {
	username string
	password string
}

func NewProfileConfigEnv() (*profileConfigEnv, error) {
	username := os.Getenv(usernameEnv)
	if username == "" {
		return nil, fmt.Errorf("environment variable %s is not set", usernameEnv)
	}

	password := os.Getenv(passwordEnv)
	if password == "" {
		return nil, fmt.Errorf("environment variable %s is not set", passwordEnv)
	}

	return &profileConfigEnv{
		username: username,
		password: password,
	}, nil
}

func (c *profileConfigEnv) Username() string {
	return c.username
}

func (c *profileConfigEnv) Password() string {
	return c.password
}
