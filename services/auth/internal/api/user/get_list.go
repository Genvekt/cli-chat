package user

import (
	"context"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
)

// GetList handles GetListRequest
func (s *Service) GetList(ctx context.Context, req *userApi.GetListRequest) (*userApi.GetListResponse, error) {
	users, err := s.userService.GetList(ctx, req.Names)
	if err != nil {
		return nil, err
	}

	return &userApi.GetListResponse{
		Users: converter.ToProtoUsersFromUsers(users),
	}, nil
}
