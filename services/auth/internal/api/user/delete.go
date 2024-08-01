package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
)

// Delete handles DeleteRequest
func (s *Service) Delete(ctx context.Context, req *userApi.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
