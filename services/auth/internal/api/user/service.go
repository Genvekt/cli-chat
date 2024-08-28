package user

import (
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils"
)

// Service is a user api implementation
type Service struct {
	userApi.UnimplementedUserV1Server
	userService service.UserService
	hasher      utils.Hasher
}

// NewService initialises user api implementation
func NewService(userService service.UserService, hasher utils.Hasher) *Service {
	return &Service{
		userService: userService,
		hasher:      hasher,
	}
}
