package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
)

// Update handles UpdateRequest
func (s *Service) Update(ctx context.Context, req *userApi.UpdateRequest) (*emptypb.Empty, error) {
	userUpdateDto := converter.ToUserUpdateDTOFromProto(req)

	err := s.userService.Update(ctx, userUpdateDto)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
