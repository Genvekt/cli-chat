package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Update applies provided function to user
func (s *userService) Update(ctx context.Context, dto *model.UserUpdateDTO) error {

	updateFunc := func(user *model.User) error {
		if dto.Name != nil {
			user.Name = *dto.Name
		}
		if dto.Email != nil {
			user.Email = *dto.Email
		}
		if dto.Role != nil {
			user.Role = *dto.Role
		}
		return nil
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepo.Update(ctx, dto.ID, updateFunc)
		if err != nil {
			return err
		}

		// Delete old user version from cache
		err = s.deleteCache(ctx, dto.ID)
		if err != nil && !errors.Is(err, ErrNoCacheUsed) {
			// We cannot leave old user version in cache when it is updated in database
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot update user: %v", err)
	}

	return nil
}
