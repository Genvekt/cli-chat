package user

import (
	"github.com/Genvekt/cli-chat/services/auth/internal/client/db"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
)

var _ repository.UserRepository = (*userService)(nil)

type userService struct {
	userRepo  repository.UserRepository
	txManager db.TxManager
}

// NewUserService initialises user service layer
func NewUserService(userRepo repository.UserRepository, txManager db.TxManager) *userService {
	return &userService{
		userRepo:  userRepo,
		txManager: txManager,
	}
}
