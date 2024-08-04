package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	chatImpl "github.com/Genvekt/cli-chat/services/chat-server/internal/api/chat"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	serviceMock "github.com/Genvekt/cli-chat/services/chat-server/internal/service/mocks"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *chatApi.SendMessageRequest
	}

	type chatServiceMockFunc func(mc minimock.MockController) *serviceMock.ChatServiceMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = fmt.Errorf("service error")

		chatId = gofakeit.Int64()
		userId = gofakeit.Int64()
		text   = gofakeit.Sentence(15)
		ts     = gofakeit.Date()

		message = &model.Message{
			SenderID:  userId,
			ChatID:    chatId,
			Content:   text,
			Timestamp: ts,
		}

		req = &chatApi.SendMessageRequest{
			Message: &chatApi.Message{
				ChatId:    chatId,
				SenderId:  userId,
				Text:      text,
				Timestamp: timestamppb.New(ts),
			},
		}

		res = &emptypb.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc minimock.MockController) *serviceMock.ChatServiceMock {
				mock := serviceMock.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "Service failure",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServiceMock: func(mc minimock.MockController) *serviceMock.ChatServiceMock {
				mock := serviceMock.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(serviceErr)
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

			apiRes, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, err)
		})
	}
}
