package user

import (
	"errors"

	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"
	"github.com/Genvekt/cli-chat/services/auth/internal/config"

	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
)

var _ service.UserService = (*userService)(nil)

var (
	// ErrNoCacheUsed indicates that cache is not initialised
	ErrNoCacheUsed = errors.New("no cache used")
)

type userService struct {
	userRepo  repository.UserRepository
	userCache repository.UserCache
	txManager db.TxManager
	config    config.UserServiceConfig
}

// NewUserService initialises user service layer
func NewUserService(
	userRepo repository.UserRepository,
	userCache repository.UserCache,
	txManager db.TxManager,
	config config.UserServiceConfig,
) *userService {
	return &userService{
		userRepo:  userRepo,
		userCache: userCache,
		txManager: txManager,
		config:    config,
	}
}
