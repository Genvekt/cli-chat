package user

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Get retrieves user by id
func (s *userService) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.getCache(ctx, id)
	if err == nil {
		return user, nil
	}

	user, err = s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user with id %d: %v", id, err)
	}

	_ = s.setCache(ctx, user)

	return user, nil
}
