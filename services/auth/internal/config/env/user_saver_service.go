package env

import (
	"fmt"
	"os"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

const (
	userSaverTopicEnv = "USER_SAVER_TOPIC"
)

var _ config.UserSaverConfig = (*userSaverConfigEnv)(nil)

type userSaverConfigEnv struct {
	topic string
}

// NewUserSaverConfigEnv reads config for user saver from env
func NewUserSaverConfigEnv() (*userSaverConfigEnv, error) {
	topic := os.Getenv(userSaverTopicEnv)
	if topic == "" {
		return nil, fmt.Errorf("environment variable %s not set", userSaverTopicEnv)
	}

	return &userSaverConfigEnv{topic: topic}, nil
}

// Topic returns kafka topic to read from
func (s *userSaverConfigEnv) Topic() string {
	return s.topic
}
