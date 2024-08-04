package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	repoMock "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/mocks"
	chatService "github.com/Genvekt/cli-chat/services/chat-server/internal/service/chat"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
	}

	type chatRepoMockFunc func(mc minimock.MockController) *repoMock.ChatRepositoryMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID  = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name             string
		args             args
		want             int64
		err              error
		chatRepoMockFunc chatRepoMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				id:  chatID,
			},
			err: nil,
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, chatID).Return(nil)
				return mock
			},
		},
		{
			name: "Chat repo error",
			args: args{
				ctx: ctx,
				id:  chatID,
			},
			err: repoErr,
			chatRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatRepositoryMock {
				mock := repoMock.NewChatRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, chatID).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatRepoMock := tt.chatRepoMockFunc(mc)
			service := chatService.NewChatService(
				chatRepoMock, nil, nil, nil, nil,
			)

			err := service.Delete(tt.args.ctx, tt.args.id)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
