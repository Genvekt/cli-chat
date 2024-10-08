package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	dbMock "github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/mocks"
	configMock "github.com/Genvekt/cli-chat/services/auth/internal/config/mocks"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestCreateWithCache(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.User
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type userCacheMockFunc func(mc minimock.MockController) *repoMock.UserCacheMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock
	type configMockFunc func(mc minimock.MockController) *configMock.UserServiceConfigMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Username()
		email    = gofakeit.Email()
		roleUser = userApi.UserRole_USER

		repoErr = fmt.Errorf("repo error")

		req = &model.User{
			Name:  name,
			Email: email,
			Role:  int(roleUser),
		}
	)

	tests := []struct {
		name              string
		args              args
		want              int64
		err               error
		userRepoMockFunc  userRepoMockFunc
		userCacheMockFunc userCacheMockFunc
		txManagerMockFunc txManagerMockFunc
		configMockFunc    configMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			userCacheMockFunc: func(mc minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "User repo create failure",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
				return mock
			},
			userCacheMockFunc: func(mc minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
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
			userRepoMock := tt.userRepoMockFunc(mc)
			userCacheMock := tt.userCacheMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)

			confMock := configMock.NewUserServiceConfigMock(mc)
			confMock.NoCacheMock.Optional().Set(func() bool { return false })

			service := userService.NewUserService(userRepoMock, userCacheMock, txManagerMock, confMock)

			servRes, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, servRes)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestCreateWithoutCache(t *testing.T) {
	type args struct {
		ctx context.Context
		req *model.User
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

		req = &model.User{
			Name:  name,
			Email: email,
			Role:  int(roleUser),
		}
	)

	tests := []struct {
		name              string
		args              args
		want              int64
		err               error
		userRepoMockFunc  userRepoMockFunc
		txManagerMockFunc txManagerMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "User repo create failure",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, repoErr)
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
			userRepoMock := tt.userRepoMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)

			confMock := configMock.NewUserServiceConfigMock(mc)
			confMock.NoCacheMock.Optional().Set(func() bool { return true })

			service := userService.NewUserService(userRepoMock, nil, txManagerMock, confMock)

			servRes, err := service.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, servRes)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
