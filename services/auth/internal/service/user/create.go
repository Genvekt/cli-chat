package user

import (
	"context"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Create creates user
func (s *userService) Create(ctx context.Context, user *model.User) (int64, error) {
	newUserID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return newUserID, nil
}
