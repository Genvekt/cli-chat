package repository

import "github.com/Genvekt/cli-chat/services/auth/model"

// UserRepository is used to manage users in some data source
type UserRepository interface {
	Get(id int64) (*model.User, error)
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) error
	Delete(id int64) error
}
