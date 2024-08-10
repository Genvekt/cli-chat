package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
)

// ToUserFromProtoInfo converts user api model to user service model
func ToUserFromProtoInfo(user *userApi.UserInfo) *model.User {
	return &model.User{
		Name:  user.Name,
		Email: user.Email,
		Role:  int(user.Role),
	}
}

// ToUserUpdateDTOFromProto packs user update requiest into user update dto for service layer
func ToUserUpdateDTOFromProto(userUpdate *userApi.UpdateRequest) *model.UserUpdateDTO {
	dto := &model.UserUpdateDTO{ID: userUpdate.Id}
	if userUpdate.GetName() != nil {
		name := userUpdate.GetName().GetValue()
		dto.Name = &name
	}
	if userUpdate.GetEmail() != nil {
		email := userUpdate.GetEmail().GetValue()
		dto.Email = &email
	}
	if userUpdate.Role != nil {
		role := int(userUpdate.GetRole())
		dto.Role = &role
	}
	return dto
}

// ToUserFiltersFromProto converts user filters from proto to service model
func ToUserFiltersFromProto(userFilters *userApi.GetListFilters) *model.UserFilters {
	return &model.UserFilters{
		Names: userFilters.GetNames(),
	}
}

// ToProtoUsersFromUsers converts slice of user service model to slice of user api model
func ToProtoUsersFromUsers(users []*model.User) []*userApi.User {
	protoUsers := make([]*userApi.User, 0, len(users))

	for _, user := range users {
		protoUsers = append(protoUsers, ToProtoUserFromUser(user))
	}

	return protoUsers
}

// ToProtoUserFromUser converts user service model to user api model
func ToProtoUserFromUser(user *model.User) *userApi.User {
	return &userApi.User{
		Id: user.ID,
		Info: &userApi.UserInfo{
			Name:  user.Name,
			Email: user.Email,
			Role:  userApi.UserRole(user.Role),
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
