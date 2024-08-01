package user

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Query retrieves users by names
func (s *userService) Query(ctx context.Context, names []string) ([]*model.User, error) {
	users, err := s.userRepo.Query(ctx, names)
	if err != nil {
		return nil, err
	}

	return users, nil
}
