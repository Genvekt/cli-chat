package auth

import (
	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
)

// Service is auth api
type Service struct {
	authApi.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewService initialises auth api implementation
func NewService(authService service.AuthService) *Service {
	return &Service{
		authService: authService,
	}
}
