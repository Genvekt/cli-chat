package repository

import (
	"context"
	"time"

	"errors"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

var ErrUserNotFound = errors.New("user not found") // ErrUserNotFound indicates user not found in data source

// UserRepository is used to manage users in some data source
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetList(ctx context.Context, filters *model.UserFilters) ([]*model.User, error)
	Update(ctx context.Context, id int64, updateFunc func(user *model.User) error) error
	Delete(ctx context.Context, id int64) error
}

// UserCache manages user cache in some data source
type UserCache interface {
	Get(ctx context.Context, id int64) (*model.User, error)
	Set(ctx context.Context, user *model.User) error
	Expire(ctx context.Context, id int64, timeout time.Duration) error
	Delete(ctx context.Context, id int64) error
}
