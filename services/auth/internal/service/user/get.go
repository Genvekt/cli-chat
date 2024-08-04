package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
)

// Get retrieves user by id
func (s *userService) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.getCache(ctx, id)
	if err == nil {
		return user, nil
	} else if !errors.Is(err, repository.ErrUserNotFound) && !errors.Is(err, ErrNoCacheUsed) {
		// We can omit cache get problems, it is not crucial for application
		// TODO: change to Error
		log.Printf("cannot get user with id %d from cache: %v", id, err)
	}

	user, err = s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user with id %d: %v", id, err)
	}

	err = s.setCache(ctx, user)
	if err != nil && !errors.Is(err, ErrNoCacheUsed) {
		// We can omit cache save problems, it is not crucial for application
		// TODO: change to Error
		log.Printf("cannot save user with id %d to cache: %v", id, err)
	}

	return user, nil
}
