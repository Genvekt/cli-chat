package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// GetList retrieves users by names
func (s *userService) GetList(ctx context.Context, names []string) ([]*model.User, error) {
	users, err := s.userRepo.GetList(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("cannot get user list: %v", err)
	}

	for _, user := range users {
		err = s.setCache(ctx, user)
		if err != nil && !errors.Is(err, ErrNoCacheUsed) {
			// We can omit cache save problems, it is not crucial for application
			// TODO: change to Error
			log.Printf("cannot save user with id %d to cache: %v", user.ID, err)
		}
	}
	return users, nil
}
