package auth

import (
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils"
)

var _ service.AuthService = (*authService)(nil)

type authService struct {
	userRepo repository.UserRepository

	refreshTokenProvider utils.TokenProvider
	accessTokenProvider  utils.TokenProvider
	hasher               utils.Hasher
}

// NewAuthService initialised auth service layer
func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenProvider utils.TokenProvider,
	accessTokenProvider utils.TokenProvider,
	hasher utils.Hasher,
) *authService {
	return &authService{
		userRepo:             userRepo,
		refreshTokenProvider: refreshTokenProvider,
		accessTokenProvider:  accessTokenProvider,
		hasher:               hasher,
	}
}
