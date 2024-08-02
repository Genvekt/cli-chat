package converter

import (
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// ToUsersFromRepo converts slice of user model of repository layer to slice of user model of service layer
func ToUsersFromRepo(users []*userApi.User) []*model.User {
	serviceUsers := make([]*model.User, 0, len(users))

	for _, user := range users {
		serviceUsers = append(serviceUsers, ToUserFromRepo(user))
	}

	return serviceUsers
}

// ToUserFromRepo converts user model of repository layer to user model of service layer
func ToUserFromRepo(user *userApi.User) *model.User {
	return &model.User{
		ID:   user.Id,
		Name: user.Info.Name,
	}
}
