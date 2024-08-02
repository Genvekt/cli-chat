package repository

import (
	"context"

	"errors"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

var ErrUserNotFound = errors.New("user not found") // ErrUserNotFound indicates user not found in data source

// UserRepository is used to manage users in some data source
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, names []string) ([]*model.User, error)
	Update(ctx context.Context, id int64, updateFunc func(user *model.User) error) error
	Delete(ctx context.Context, id int64) error
}
