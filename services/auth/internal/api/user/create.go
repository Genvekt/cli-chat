package user

import (
	"context"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
)

// Create handles CreateRequest
func (s *Service) Create(ctx context.Context, req *userApi.CreateRequest) (*userApi.CreateResponse, error) {
	userID, err := s.userService.Create(ctx, converter.ToUserFromProtoInfo(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	return &userApi.CreateResponse{
		Id: userID,
	}, nil
}
