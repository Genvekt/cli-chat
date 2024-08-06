package service

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// UserService is a service layer with business logic related to users
type UserService interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, names []string) ([]*model.User, error)
	Update(ctx context.Context, id int64, updateFunc func(user *model.User) error) error
	Delete(ctx context.Context, id int64) error
}
