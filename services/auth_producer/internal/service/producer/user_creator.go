package producer

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	gofakeit "github.com/brianvoe/gofakeit/v7"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/service"
)

var _ service.UserCreatorService = (*userCreatorService)(nil)

type userCreatorService struct {
	config   config.UserCreatorConfig
	producer kafka.Producer
	messages chan *sarama.ProducerMessage
}

// NewUserCreatorService initialises user creator instance
func NewUserCreatorService(conf config.UserCreatorConfig, producer kafka.Producer) *userCreatorService {
	return &userCreatorService{
		config:   conf,
		producer: producer,
		messages: make(chan *sarama.ProducerMessage),
	}
}

// Create creates random user and sends it into kafka
func (s *userCreatorService) Create(_ context.Context) error {
	user := &userApi.UserInfo{
		Name:  gofakeit.Username(),
		Email: gofakeit.Email(),
	}

	encodedUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: s.config.Topic(),
		Value: sarama.StringEncoder(encodedUser),
	}

	s.messages <- msg

	return nil
}

// RunProducer starts producer loop
func (s *userCreatorService) RunProducer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-s.messages:
			err := s.produce(ctx, msg)
			if err != nil {
				return err
			}
		}
	}
}

func (s *userCreatorService) produce(_ context.Context, msg *sarama.ProducerMessage) error {
	err := s.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
