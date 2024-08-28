package kafka

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	kafkaCli "github.com/Genvekt/cli-chat/services/auth_producer/internal/client/kafka"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/client/kafka/user/converter"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/model"
)

var _ kafkaCli.UserClient = (*userClient)(nil)

type userClient struct {
	producer kafka.Producer[sarama.ProducerMessage]
	config   config.UserKafkaClientConfig
}

// NewUserClient initialises kafka user client
func NewUserClient(config config.UserKafkaClientConfig, producer kafka.Producer[sarama.ProducerMessage]) *userClient {
	return &userClient{
		config:   config,
		producer: producer,
	}
}

// Create creates user
func (c *userClient) Create(_ context.Context, user *model.User) error {
	req := converter.UserToKafkaUser(user)

	msg, err := converter.StructToMsg(req, c.config.Topic())
	if err != nil {
		return err
	}

	err = c.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
