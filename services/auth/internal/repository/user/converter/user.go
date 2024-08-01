package converter

import (
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/model"
)

// ToUserFromRepo converts user repository model to user service model
func ToUserFromRepo(user *repoModel.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
