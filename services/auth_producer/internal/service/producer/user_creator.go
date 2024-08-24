package producer

import (
	"context"

	gofakeit "github.com/brianvoe/gofakeit/v7"

	kafkaCli "github.com/Genvekt/cli-chat/services/auth_producer/internal/client/kafka"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/model"

	"github.com/Genvekt/cli-chat/services/auth_producer/internal/service"
)

var _ service.UserCreatorService = (*userCreatorService)(nil)

type userCreatorService struct {
	userClient kafkaCli.UserClient
}

// NewUserCreatorService initialises user creator instance
func NewUserCreatorService(userClient kafkaCli.UserClient) *userCreatorService {
	return &userCreatorService{
		userClient: userClient,
	}
}

// Create creates random user and sends it into kafka
func (s *userCreatorService) Create(ctx context.Context) error {
	user := &model.User{
		Name:  gofakeit.Username(),
		Email: gofakeit.Email(),
	}

	err := s.userClient.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
