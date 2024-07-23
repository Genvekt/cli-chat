package repository

import "github.com/Genvekt/cli-chat/services/auth/model"

type UserRepository interface {
	Get(id int64) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) error
	Delete(id int64) error
}
