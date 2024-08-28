package access

import (
	accessApi "github.com/Genvekt/cli-chat/libraries/api/access/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
)

// Service is access api
type Service struct {
	accessApi.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewService initialises access api implementation
func NewService(accessService service.AccessService) *Service {
	return &Service{
		accessService: accessService,
	}
}
