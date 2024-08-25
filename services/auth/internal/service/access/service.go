package access

import (
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils"
)

type accessService struct {
	accessTokenProvider utils.TokenProvider
	accessRepo          repository.AccessRepository
}

// NewAccessService initialised access service layer
func NewAccessService(accessTokenProvider utils.TokenProvider, accessRepo repository.AccessRepository) *accessService {
	return &accessService{
		accessTokenProvider: accessTokenProvider,
		accessRepo:          accessRepo,
	}
}
