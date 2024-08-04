package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/client/db"
	dbMock "github.com/Genvekt/cli-chat/services/chat-server/internal/client/db/mocks"
	serviceMock "github.com/Genvekt/cli-chat/services/chat-server/internal/client/service/mocks"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"

	repoMock "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/mocks"
	chatService "github.com/Genvekt/cli-chat/services/chat-server/internal/service/chat"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		name      string
		usernames []string
	}

	type authClientMockFunc func(mc minimock.MockController) *serviceMock.AuthClientMock
	type chatRepoMockFunc func(mc minimock.MockController) *repoMock.ChatRepositoryMock
	type chatMemberRepoMockFunc func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID   = gofakeit.Int64()
		name     = gofakeit.Username()
		userID   = gofakeit.Int64()
		username = gofakeit.Username()
		user     = &model.User{ID: userID, Name: username}

		userIDs   = []int64{userID}
		usernames = []string{username}
		users     = []*model.User{user}

		chat = &model.Chat{
			Name: name,
		}

		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name                   string
		args                   args
		want                   int64
		err                    error
		authClientMockFunc     authClientMockFunc
		chatRepoMockFunc       chatRepoMockFunc
		chatMemberRepoMockFunc chatMemberRepoMockFunc
		txManagerMockFunc      txManagerMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       ctx,
				name:      name,
				usernames: usernames,
			},
			want: chatID,
			err:  nil,
			authClientMockFunc: func(mc minimock.MockController) *serviceMock.AuthClientMock {
				mock := serviceMock.NewAuthClientMock(mc)
				mock.GetListMock.Expect(ctx, usernames).Return(users, nil)
				return mock
			},
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chat).Return(chatID, nil)
				return mock
			},
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.CreateBatchMock.Expect(ctx, chatID, userIDs).Return(nil)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				mock := dbMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Auth client empty response",
			args: args{
				ctx:       ctx,
				name:      name,
				usernames: usernames,
			},
			want: 0,
			err:  repoErr,
			authClientMockFunc: func(mc minimock.MockController) *serviceMock.AuthClientMock {
				mock := serviceMock.NewAuthClientMock(mc)
				mock.GetListMock.Expect(ctx, usernames).Return(nil, nil)
				return mock
			},
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				return mock
			},
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				mock := dbMock.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "Auth client error",
			args: args{
				ctx:       ctx,
				name:      name,
				usernames: usernames,
			},
			want: 0,
			err:  repoErr,
			authClientMockFunc: func(mc minimock.MockController) *serviceMock.AuthClientMock {
				mock := serviceMock.NewAuthClientMock(mc)
				mock.GetListMock.Expect(ctx, usernames).Return(nil, repoErr)
				return mock
			},
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				return mock
			},
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				mock := dbMock.NewTxManagerMock(mc)
				return mock
			},
		},
		{
			name: "Chat repo error",
			args: args{
				ctx:       ctx,
				name:      name,
				usernames: usernames,
			},
			want: 0,
			err:  repoErr,
			authClientMockFunc: func(mc minimock.MockController) *serviceMock.AuthClientMock {
				mock := serviceMock.NewAuthClientMock(mc)
				mock.GetListMock.Expect(ctx, usernames).Return(users, nil)
				return mock
			},
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chat).Return(0, repoErr)
				return mock
			},
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				mock := dbMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Chat member repo error",
			args: args{
				ctx:       ctx,
				name:      name,
				usernames: usernames,
			},
			want: 0,
			err:  repoErr,
			authClientMockFunc: func(mc minimock.MockController) *serviceMock.AuthClientMock {
				mock := serviceMock.NewAuthClientMock(mc)
				mock.GetListMock.Expect(ctx, usernames).Return(users, nil)
				return mock
			},
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chat).Return(chatID, nil)
				return mock
			},
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.CreateBatchMock.Expect(ctx, chatID, userIDs).Return(repoErr)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				mock := dbMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "Tx manager error",
			args: args{
				ctx:       ctx,
				name:      name,
				usernames: usernames,
			},
			want: 0,
			err:  repoErr,
			authClientMockFunc: func(mc minimock.MockController) *serviceMock.AuthClientMock {
				mock := serviceMock.NewAuthClientMock(mc)
				mock.GetListMock.Expect(ctx, usernames).Return(users, nil)
				return mock
			},
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chat).Return(chatID, nil)
				return mock
			},
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.CreateBatchMock.Expect(ctx, chatID, userIDs).Return(nil)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				mock := dbMock.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, f db.Handler) (err error) {
					_ = f(ctx)
					return repoErr

				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authClientMock := tt.authClientMockFunc(mc)
			chatRepoMock := tt.chatRepoMockFunc(mc)
			chatMemberRepoMock := tt.chatMemberRepoMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)
			service := chatService.NewChatService(
				chatRepoMock,
				chatMemberRepoMock,
				nil,
				authClientMock,
				txManagerMock,
			)

			servRes, err := service.Create(tt.args.ctx, tt.args.name, tt.args.usernames)
			require.Equal(t, tt.want, servRes)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
