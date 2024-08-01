package user

import (
	"context"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
)

// Query handles QueryRequest
func (s *Service) Query(ctx context.Context, req *userApi.QueryRequest) (*userApi.QueryResponse, error) {
	users, err := s.userService.Query(ctx, req.Names)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*userApi.User, 0, len(users))
	for _, user := range users {
		protoUsers = append(protoUsers, converter.ToProtoUserFromUser(user))
	}

	return &userApi.QueryResponse{
		Users: protoUsers,
	}, nil
}
