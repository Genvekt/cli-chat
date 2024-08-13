package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	"github.com/Genvekt/cli-chat/services/auth/internal/config"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
)

var _ service.ConsumerService = (*userSaverService)(nil)

type userSaverService struct {
	config   config.UserSaverConfig
	consumer kafka.Consumer
	userRepo repository.UserRepository
}

// NewUserSaverService inits instance of user saver
func NewUserSaverService(
	conf config.UserSaverConfig,
	consumer kafka.Consumer,
	userRepo repository.UserRepository,
) *userSaverService {
	return &userSaverService{
		config:   conf,
		consumer: consumer,
		userRepo: userRepo,
	}
}

// RunConsumer starts consumer loop
func (s *userSaverService) RunConsumer(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-s.run(ctx):
			if err != nil {
				return err
			}
		}
	}
}

func (s *userSaverService) run(ctx context.Context) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		errChan <- s.consumer.Consume(ctx, s.config.Topic(), s.UserSaveHandler)
	}()

	return errChan
}

// UserSaveHandler processes msg from kafka
func (s *userSaverService) UserSaveHandler(ctx context.Context, iMsg interface{}) error {
	msg, ok := iMsg.(*sarama.ConsumerMessage)
	if !ok {
		return fmt.Errorf("expected a ConsumerMessage")
	}

	userInfo := &userApi.UserInfo{}
	err := json.Unmarshal(msg.Value, userInfo)
	if err != nil {
		return err
	}

	user := converter.ToUserFromProtoInfo(userInfo)

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	log.Printf("user with id %d created from kafka\n", id)

	return nil
}
