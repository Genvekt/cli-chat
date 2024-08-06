package user

import (
	"context"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
)

// Get handles GetRequest
func (s *Service) Get(ctx context.Context, req *userApi.GetRequest) (*userApi.GetResponse, error) {
	user, err := s.userService.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userApi.GetResponse{
		User: converter.ToProtoUserFromUser(user),
	}, nil
}
