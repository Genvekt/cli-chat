package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/model"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	repoMock "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/mocks"
	chatService "github.com/Genvekt/cli-chat/services/chat-server/internal/service/chat"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		message *model.Message
	}

	type chatMemberRepoMockFunc func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock
	type messageRepoMockFunc func(mc minimock.MockController) *repoMock.MessageRepositoryMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		chatID = gofakeit.Int64()
		userID = gofakeit.Int64()

		message = &model.Message{
			ChatID:    chatID,
			SenderID:  userID,
			Content:   gofakeit.Sentence(5),
			Timestamp: gofakeit.Date(),
		}

		chatMember = &model.ChatMember{
			ID:       userID,
			JoinedAt: gofakeit.Date(),
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name                   string
		args                   args
		err                    error
		messageRepoMockFunc    messageRepoMockFunc
		chatMemberRepoMockFunc chatMemberRepoMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: nil,
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.GetMock.Expect(ctx, chatID, userID).Return(chatMember, nil)
				return mock
			},
			messageRepoMockFunc: func(mc minimock.MockController) *repoMock.MessageRepositoryMock {
				mock := repoMock.NewMessageRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "Chat member not found",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: repoErr,
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.GetMock.Expect(ctx, chatID, userID).Return(nil, repository.ErrChatMemberNotFound)
				return mock
			},
			messageRepoMockFunc: func(mc minimock.MockController) *repoMock.MessageRepositoryMock {
				mock := repoMock.NewMessageRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "Chat member repo error",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: repoErr,
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.GetMock.Expect(ctx, chatID, userID).Return(nil, repoErr)
				return mock
			},
			messageRepoMockFunc: func(mc minimock.MockController) *repoMock.MessageRepositoryMock {
				mock := repoMock.NewMessageRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "Message repo error",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: repoErr,
			chatMemberRepoMockFunc: func(mc minimock.MockController) *repoMock.ChatMemberRepositoryMock {
				mock := repoMock.NewChatMemberRepositoryMock(mc)
				mock.GetMock.Expect(ctx, chatID, userID).Return(chatMember, nil)
				return mock
			},
			messageRepoMockFunc: func(mc minimock.MockController) *repoMock.MessageRepositoryMock {
				mock := repoMock.NewMessageRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, message).Return(repoErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			chatMemberRepoMock := tt.chatMemberRepoMockFunc(mc)
			messageRepoMock := tt.messageRepoMockFunc(mc)
			service := chatService.NewChatService(
				nil, chatMemberRepoMock, messageRepoMock, nil, nil,
			)

			err := service.SendMessage(tt.args.ctx, tt.args.message)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
