package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	dbMock "github.com/Genvekt/cli-chat/services/auth/internal/client/db/mocks"
	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req int64
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name             string
		args             args
		err              error
		userRepoMockFunc userRepoMockFunc
		txManagerMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
				req: id,
			},
			err: repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
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

			err := service.Delete(tt.args.ctx, tt.args.req)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
