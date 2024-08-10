package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	dbMock "github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/mocks"
	configMock "github.com/Genvekt/cli-chat/services/auth/internal/config/mocks"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestGetWithCache(t *testing.T) {
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

		id       = gofakeit.Int64()
		name     = gofakeit.Username()
		email    = gofakeit.Email()
		roleUser = userApi.UserRole_USER

		cacheTTL = time.Duration(gofakeit.Minute()) * time.Minute

		repoErr = fmt.Errorf("repo error")

		user = &model.User{
			ID:    id,
			Name:  name,
			Email: email,
			Role:  int(roleUser),
		}
	)

	tests := []struct {
		name              string
		args              args
		want              *model.User
		err               error
		userCacheMockFunc userCacheMockFunc
		userRepoMockFunc  userRepoMockFunc
		configMockFunc    configMockFunc
		txManagerMockFunc txManagerMockFunc
	}{
		{
			name: "Success from cache",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: user,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				mock.ExpireMock.Expect(ctx, id, cacheTTL).Return(nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.CacheTTLMock.Return(cacheTTL)
				mock.NoCacheMock.Return(false)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "Success from DB",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: user,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, repository.ErrUserNotFound)
				mock.SetMock.Expect(ctx, user).Return(nil)
				mock.ExpireMock.Expect(ctx, id, cacheTTL).Return(nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.CacheTTLMock.Return(cacheTTL)
				mock.NoCacheMock.Return(false)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "Cache get failure",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: user,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				mock.SetMock.Expect(ctx, user).Return(nil)
				mock.ExpireMock.Expect(ctx, id, cacheTTL).Return(nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.CacheTTLMock.Return(cacheTTL)
				mock.NoCacheMock.Return(false)
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
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(false)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repository.ErrUserNotFound)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "User cache set failure",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: user,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(false)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repository.ErrUserNotFound)
				mock.SetMock.Expect(ctx, user).Return(repoErr)
				return mock
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "User cache set ttl failure",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: user,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.CacheTTLMock.Return(cacheTTL)
				mock.NoCacheMock.Return(false)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repository.ErrUserNotFound)
				mock.SetMock.Expect(ctx, user).Return(nil)
				mock.ExpireMock.Expect(ctx, id, cacheTTL).Return(repoErr)
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
			confMock := tt.configMockFunc(mc)
			service := userService.NewUserService(userRepoMock, userCacheMock, txManagerMock, confMock)

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

func TestGetWithoutCache(t *testing.T) {
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

		id       = gofakeit.Int64()
		name     = gofakeit.Username()
		email    = gofakeit.Email()
		roleUser = userApi.UserRole_USER

		repoErr = fmt.Errorf("repo error")

		user = &model.User{
			ID:    id,
			Name:  name,
			Email: email,
			Role:  int(roleUser),
		}
	)

	tests := []struct {
		name              string
		args              args
		want              *model.User
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
			want: user,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(user, nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(true)
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
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(true)
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
			confMock := tt.configMockFunc(mc)
			service := userService.NewUserService(userRepoMock, nil, txManagerMock, confMock)

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
