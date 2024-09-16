package producer

import (
	"context"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"go.uber.org/zap"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
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
		Name:     gofakeit.Username(),
		Email:    gofakeit.Email(),
		Role:     gofakeit.RandomInt([]int{int(userApi.UserRole_USER), int(userApi.UserRole_ADMIN)}),
		Password: gofakeit.Password(true, true, true, false, false, 20),
	}

	err := s.userClient.Create(ctx, user)
	if err != nil {
		return err
	}

	logger.Info("Created user %s with password %s",
		zap.String("username", user.Name),
		zap.String("password", user.Password),
	)

	return nil
}
