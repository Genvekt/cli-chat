package service

import (
	"context"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// AuthClient abstraction of connection to auth service
type AuthClient interface {
	GetList(ctx context.Context, usernames []string) ([]*model.User, error)
}

// UserGrpcClient wrapper around grpc client to auth service
type UserGrpcClient interface {
	GetList(ctx context.Context, req *userApi.GetListRequest) (*userApi.GetListResponse, error)
}
