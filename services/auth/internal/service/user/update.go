package user

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Update applies provided function to user
func (s *userService) Update(ctx context.Context, id int64, updateFunc func(user *model.User) error) error {
	err := s.userRepo.Update(ctx, id, updateFunc)
	if err != nil {
		return err
	}

	return nil
}
