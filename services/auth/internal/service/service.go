package service

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// UserService is a service layer with business logic related to users
type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, filters *model.UserFilters) ([]*model.User, error)
	Update(ctx context.Context, dto *model.UserUpdateDTO) error
	Delete(ctx context.Context, id int64) error
}

// ConsumerService consumes messages
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
