package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// Update handles UpdateRequest
func (s *Service) Update(ctx context.Context, req *userApi.UpdateRequest) (*emptypb.Empty, error) {
	updateFunc := func(user *model.User) error {
		if req.Email != nil {
			user.Email = req.Email.Value
		}
		if req.Name != nil {
			user.Name = req.Name.Value
		}
		if req.Role != nil {
			user.Role = int(*req.Role)
		}
		return nil
	}

	err := s.userService.Update(ctx, req.Id, updateFunc)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
