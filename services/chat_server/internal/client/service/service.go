package service

import (
	"context"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
)

// AuthClient abstraction of connection to auth service
type AuthClient interface {
	Query(ctx context.Context, userbames []string) ([]*model.User, error)
	Close() error
}

// UserGrpcClient wrapper around grpc client to auth service
type UserGrpcClient interface {
	Query(ctx context.Context, req *userApi.QueryRequest) (*userApi.QueryResponse, error)
}
