package repository

import (
	"context"
	"time"

	"errors"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

var ErrUserNotFound = errors.New("user not found")        // ErrUserNotFound indicates user not found in data source
var ErrRuleNotFound = errors.New("access rule not found") // ErrRuleNotFound indicates access rule not found in data source

// UserRepository is used to manage users in some data source
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetList(ctx context.Context, filters *model.UserFilters) ([]*model.User, error)
	Update(ctx context.Context, id int64, updateFunc func(user *model.User) error) error
	Delete(ctx context.Context, id int64) error
}

// AccessRepository is used to manage access controls in some data source
type AccessRepository interface {
	GetEndpointAccessRule(ctx context.Context, endpoint string) (*model.EndpointAccessRule, error)
}

// UserCache manages user cache in some data source
type UserCache interface {
	Get(ctx context.Context, id int64) (*model.User, error)
	Set(ctx context.Context, user *model.User) error
	Expire(ctx context.Context, id int64, timeout time.Duration) error
	Delete(ctx context.Context, id int64) error
}
