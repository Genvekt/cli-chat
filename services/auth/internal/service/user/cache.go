package user

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// setCache saves user into cache and sets ttl
func (s *userService) setCache(ctx context.Context, user *model.User) error {
	if s.config.NoCache() {
		return nil
	}

	err := s.userCache.Set(ctx, user)
	if err != nil {
		return err
	}

	// Set timeout for user
	err = s.userCache.Expire(ctx, user.ID, s.config.CacheTTL())
	if err != nil {
		return err
	}

	return nil
}

// getCache retrieves user from cache by id and resets ttl
func (s *userService) getCache(ctx context.Context, id int64) (*model.User, error) {
	if s.config.NoCache() {
		return nil, ErrNoCacheUsed
	}

	user, err := s.userCache.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Reset timeout for user as it was recently retrieved and may be retrieved once more
	err = s.userCache.Expire(ctx, id, s.config.CacheTTL())
	if err != nil {
		return nil, err
	}

	return user, nil
}

// deleteCache removes user from cache
func (s *userService) deleteCache(ctx context.Context, id int64) error {
	if s.config.NoCache() {
		return nil
	}

	err := s.userCache.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
