package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	accessApi "github.com/Genvekt/cli-chat/libraries/api/access/v1"
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// AuthClient abstraction of connection to auth service
type AuthClient interface {
	GetList(ctx context.Context, usernames []string) ([]*model.User, error)
}

// AccessClient abstraction of connection to access api
type AccessClient interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

// UserGrpcClient wrapper around grpc client to auth service
type UserGrpcClient interface {
	GetList(ctx context.Context, req *userApi.GetListRequest) (*userApi.GetListResponse, error)
}

// AccessGrpcClient wrapper around grpc client to access service
type AccessGrpcClient interface {
	Check(ctx context.Context, req *accessApi.CheckRequest) (*emptypb.Empty, error)
}
