package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	chatImpl "github.com/Genvekt/cli-chat/services/chat-server/internal/api/chat"
	serviceMock "github.com/Genvekt/cli-chat/services/chat-server/internal/service/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *chatApi.CreateRequest
	}

	type chatServiceMockFunc func(mc minimock.MockController) *serviceMock.ChatServiceMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		chatName  = gofakeit.Sentence(3)
		usernames = []string{gofakeit.Username()}

		serviceErr = fmt.Errorf("service error")

		reqWithUsers = &chatApi.CreateRequest{
			Name:      chatName,
			Usernames: usernames,
		}

		reqNoUsers = &chatApi.CreateRequest{
			Name:      chatName,
			Usernames: nil,
		}

		res = &chatApi.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *chatApi.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: reqWithUsers,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc minimock.MockController) *serviceMock.ChatServiceMock {
				mock := serviceMock.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctx, chatName, usernames).Return(id, nil)
				return mock
			},
		},
		{
			name: "No users",
			args: args{
				ctx: ctx,
				req: reqNoUsers,
			},
			want: nil,
			err:  chatImpl.ErrEmptyChat,
			chatServiceMock: func(mc minimock.MockController) *serviceMock.ChatServiceMock {
				mock := serviceMock.NewChatServiceMock(mc)
				return mock
			},
		},
		{
			name: "Service failure",
			args: args{
				ctx: ctx,
				req: reqWithUsers,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc minimock.MockController) *serviceMock.ChatServiceMock {
				mock := serviceMock.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctx, chatName, usernames).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.chatServiceMock(mc)
			api := chatImpl.NewService(userServiceMock)

			apiRes, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, err)
		})
	}
}
