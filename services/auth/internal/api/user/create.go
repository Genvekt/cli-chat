package user

import (
	"context"
	"errors"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/converter"
)

// ErrPasswordsNotMatch indicates that password is not equal to password confirm
var ErrPasswordsNotMatch = errors.New(`passwords do not match`)

func validateCreateRequest(req *userApi.CreateRequest) error {
	// proto validate cannot compare fields..
	if req.Password != req.PasswordConfirm {
		return ErrPasswordsNotMatch
	}

	return nil
}

// Create handles CreateRequest
func (s *Service) Create(ctx context.Context, req *userApi.CreateRequest) (*userApi.CreateResponse, error) {
	err := validateCreateRequest(req)
	if err != nil {
		return nil, err
	}

	passwordHash, err := s.hasher.HashPassword(ctx, req.Password)
	if err != nil {
		return nil, err
	}

	userID, err := s.userService.Create(ctx, converter.ToUserFromProtoInfo(req.GetInfo(), passwordHash))
	if err != nil {
		return nil, err
	}

	return &userApi.CreateResponse{
		Id: userID,
	}, nil
}
