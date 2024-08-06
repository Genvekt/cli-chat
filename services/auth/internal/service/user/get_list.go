package user

import (
	"context"
	"fmt"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// GetList retrieves users by names
func (s *userService) GetList(ctx context.Context, names []string) ([]*model.User, error) {
	users, err := s.userRepo.GetList(ctx, names)
	if err != nil {
		return nil, fmt.Errorf("cannot get user list: %v", err)
	}

	return users, nil
}
