package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	dbMock "github.com/Genvekt/cli-chat/services/auth/internal/client/db/mocks"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name             string
		args             args
		want             *model.User
		err              error
		userRepoMockFunc userRepoMockFunc
		txManagerMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Return(nil)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "User repo failure",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Return(repoErr)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepoMock := tt.userRepoMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)
			service := userService.NewUserService(userRepoMock, txManagerMock)

			err := service.Update(tt.args.ctx, tt.args.id, func(user *model.User) error { return nil })
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
