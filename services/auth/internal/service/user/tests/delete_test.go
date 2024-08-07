package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"
	dbMock "github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/mocks"
	configMock "github.com/Genvekt/cli-chat/services/auth/internal/config/mocks"

	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestDeleteWithCache(t *testing.T) {
	type args struct {
		ctx context.Context
		req int64
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type userCacheMockFunc func(mc minimock.MockController) *repoMock.UserCacheMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock
	type configMockFunc func(mc minimock.MockController) *configMock.UserServiceConfigMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name              string
		args              args
		err               error
		userCacheMockFunc userCacheMockFunc
		userRepoMockFunc  userRepoMockFunc
		configMockFunc    configMockFunc
		txManagerMockFunc txManagerMockFunc
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
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.ExpireMock.Expect(ctx, id, 0).Return(nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
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
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				return nil
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
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
			name: "User cache failure",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.ExpireMock.Expect(ctx, id, 0).Return(repoErr)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := tt.userRepoMockFunc(mc)
			userCacheMock := tt.userCacheMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)
			confMock := tt.configMockFunc(mc)
			service := userService.NewUserService(userRepoMock, userCacheMock, txManagerMock, confMock)

			err := service.Delete(tt.args.ctx, tt.args.req)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestDeleteWithoutCache(t *testing.T) {
	type args struct {
		ctx context.Context
		req int64
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock
	type configMockFunc func(mc minimock.MockController) *configMock.UserServiceConfigMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name              string
		args              args
		err               error
		userRepoMockFunc  userRepoMockFunc
		configMockFunc    configMockFunc
		txManagerMockFunc txManagerMockFunc
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
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
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
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := tt.userRepoMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)
			confMock := tt.configMockFunc(mc)
			service := userService.NewUserService(userRepoMock, nil, txManagerMock, confMock)

			err := service.Delete(tt.args.ctx, tt.args.req)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
