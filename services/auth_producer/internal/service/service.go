package service

import "context"

// UserCreatorService produces new users
type UserCreatorService interface {
	Create(ctx context.Context) error
}
