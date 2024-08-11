package converter

import (
	"time"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/redis/model"
)

// ToUserFromRepo converts user repository model to user service model
func ToUserFromRepo(user *repoModel.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: time.Unix(0, user.UpdatedAtNs),
	}
}

// ToRepoFromUser converts user service model to user repository model
func ToRepoFromUser(user *model.User) *repoModel.User {
	return &repoModel.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		CreatedAtNs: user.CreatedAt.UnixNano(),
		UpdatedAtNs: user.UpdatedAt.UnixNano(),
	}
}
