package repository

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/model"
)

// UserRepository is used to manage users in some data source
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, updateFunc func(user *model.User) error) error
	Delete(ctx context.Context, id int64) error
}
