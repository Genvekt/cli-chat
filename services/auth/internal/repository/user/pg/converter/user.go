package converter

import (
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoModel "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/pg/model"
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

// ToUsersFromRepo converts slice of user repository model to slice of user service model
func ToUsersFromRepo(users []*repoModel.User) []*model.User {
	serviceUsers := make([]*model.User, 0, len(users))
	for _, user := range users {
		serviceUsers = append(serviceUsers, ToUserFromRepo(user))
	}
	return serviceUsers
}
