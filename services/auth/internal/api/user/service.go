package user

import (
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
)

// Service is a user api implementation
type Service struct {
	userApi.UnimplementedUserV1Server
	userService service.UserService
}

// NewService initialises user api implementation
func NewService(userService service.UserService) *Service {
	return &Service{
		userService: userService,
	}
}
