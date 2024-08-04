package user

import (
	"errors"
	"reflect"
	"time"

	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"

	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
)

var _ service.UserService = (*userService)(nil)

var (
	// ErrNoCacheUsed indicates that cache is not initialised
	ErrNoCacheUsed = errors.New("no cache used")
	cacheTimeout   = time.Minute
)

type userService struct {
	userRepo  repository.UserRepository
	userCache repository.UserCache
	txManager db.TxManager
}

// NewUserService initialises user service layer
func NewUserService(
	userRepo repository.UserRepository,
	userCache repository.UserCache,
	txManager db.TxManager,
) *userService {
	return &userService{
		userRepo:  userRepo,
		userCache: userCache,
		txManager: txManager,
	}
}

func (s *userService) isCacheUsed() bool {
	return !(s.userCache == nil || reflect.ValueOf(s.userCache).IsNil())
}
