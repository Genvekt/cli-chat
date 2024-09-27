package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/model"
)

type ChatClient interface {
	Create(ctx context.Context, name string, usernames []string) (int64, error)
	SendMessage(ctx context.Context, message *model.Message) error
	Connect(ctx context.Context, chatID int64, username string) (chan *model.Message, error)
}

type ChatGRPCClient interface {
	Create(ctx context.Context, req *chatApi.CreateRequest) (*chatApi.CreateResponse, error)
	SendMessage(ctx context.Context, req *chatApi.SendMessageRequest) (*emptypb.Empty, error)
	Connect(ctx context.Context, req *chatApi.ConnectChatRequest) (chatApi.ChatV1_ConnectChatClient, error)
}

type AuthClient interface {
	Login(ctx context.Context, username string, password string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AuthGRPCClient interface {
	Login(ctx context.Context, req *authApi.LoginRequest) (*authApi.LoginResponse, error)
	GetAccessToken(ctx context.Context, req *authApi.GetAccessTokenRequest) (*authApi.GetAccessTokenResponse, error)
}
