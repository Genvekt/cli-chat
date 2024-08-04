package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	dbMock "github.com/Genvekt/cli-chat/services/auth/internal/client/db/mocks"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestGet(t *testing.T) {
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

		id       = gofakeit.Int64()
		name     = gofakeit.Username()
		email    = gofakeit.Email()
		roleUser = userApi.UserRole_USER

		repoErr = fmt.Errorf("repo error")

		res = &model.User{
			ID:    id,
			Name:  name,
			Email: email,
			Role:  int(roleUser),
		}
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
				req: id,
			},
			want: res,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
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
			want: nil,
			err:  repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
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

			servRes, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, servRes)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
