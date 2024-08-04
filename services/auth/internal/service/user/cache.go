package user

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// setCache saves user into cache and sets ttl
func (s *userService) setCache(ctx context.Context, user *model.User) error {
	if !s.isCacheUsed() {
		return ErrNoCacheUsed
	}

	err := s.userCache.Set(ctx, user)
	if err != nil {
		return err
	}

	// Set timeout for user
	err = s.userCache.Expire(ctx, user.ID, cacheTimeout)
	if err != nil {
		return err
	}

	return nil
}

// getCache retrieves user from cache by id and resets ttl
func (s *userService) getCache(ctx context.Context, id int64) (*model.User, error) {
	if !s.isCacheUsed() {
		return nil, ErrNoCacheUsed
	}

	user, err := s.userCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Reset timeout for user as it was recently retrieved
	err = s.userCache.Expire(ctx, id, cacheTimeout)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// deleteCache removes user from cache
func (s *userService) deleteCache(ctx context.Context, id int64) error {
	if !s.isCacheUsed() {
		return ErrNoCacheUsed
	}

	err := s.userCache.Expire(ctx, id, 0)
	if err != nil {
		return err
	}

	return nil
}
