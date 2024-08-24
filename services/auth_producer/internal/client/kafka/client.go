package kafka

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth_producer/internal/model"
)

// UserClient client that manages users somewhere
type UserClient interface {
	Create(ctx context.Context, user *model.User) error
}
