package user

import (
	"context"
	"errors"
	"fmt"
)

// Delete deletes user by id
func (s *userService) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		txErr := s.userRepo.Delete(ctx, id)
		if txErr != nil {
			return txErr
		}

		txErr = s.deleteCache(ctx, id)
		if txErr != nil && !errors.Is(txErr, ErrNoCacheUsed) {
			// We cannot leave user in cache when it is deleted from database
			return txErr
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot delete user with id %d from cache: %v", id, err)
	}

	return nil
}
