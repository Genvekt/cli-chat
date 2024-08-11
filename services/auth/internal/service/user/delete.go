package user

import (
	"context"
	"errors"
	"fmt"
)

// Delete deletes user by id
func (s *userService) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.userRepo.Delete(ctx, id)
		if err != nil {
			return err
		}

		err = s.deleteCache(ctx, id)
		if err != nil && !errors.Is(err, ErrNoCacheUsed) {
			// We cannot leave user in cache when it is deleted from database
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot delete user with id %d from cache: %v", id, err)
	}

	return nil
}
