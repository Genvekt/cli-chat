package env

import (
	"fmt"
	"os"

	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
)

const (
	userCreatorTopicEnv = "USER_CREATOR_TOPIC"
)

var _ config.UserKafkaClientConfig = (*userCreatorConfigEnv)(nil)

type userCreatorConfigEnv struct {
	topic string
}

// NewUserCreatorConfigEnv reads env config for user creator service
func NewUserCreatorConfigEnv() (*userCreatorConfigEnv, error) {
	topic := os.Getenv(userCreatorTopicEnv)
	if topic == "" {
		return nil, fmt.Errorf("environment variable %s not set", userCreatorTopicEnv)
	}

	return &userCreatorConfigEnv{topic: topic}, nil
}

// Topic returns tha name of topic to publish to
func (s *userCreatorConfigEnv) Topic() string {
	return s.topic
}
