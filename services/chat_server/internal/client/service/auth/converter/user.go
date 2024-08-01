package converter

import (
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// ToUserFromRepo forms user model of service layer from user model of repository layer
func ToUserFromRepo(user *userApi.User) *model.User {
	return &model.User{
		ID:   user.Id,
		Name: user.Info.Name,
	}
}
