package user

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// GetList retrieves users by names
func (s *userService) GetList(ctx context.Context, filters *model.UserFilters) ([]*model.User, error) {
	users, err := s.userRepo.GetList(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("cannot get user list: %v", err)
	}

	for _, user := range users {
		_ = s.setCache(ctx, user)
	}
	return users, nil
}
