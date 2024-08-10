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

func TestGetListWithoutCache(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *model.UserFilters
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type userCacheMockFunc func(mc minimock.MockController) *repoMock.UserCacheMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name      = gofakeit.Username()
		usernames = []string{gofakeit.Username(), name}

		req = &model.UserFilters{
			Names: usernames,
		}

		repoErr = fmt.Errorf("repo error")

		res = []*model.User{{
			ID:    gofakeit.Int64(),
			Name:  name,
			Email: gofakeit.Email(),
			Role:  int(userApi.UserRole_USER),
		}}
	)

	tests := []struct {
		name              string
		args              args
		want              []*model.User
		err               error
		userCacheMockFunc userCacheMockFunc
		userRepoMockFunc  userRepoMockFunc
		txManagerMockFunc txManagerMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetListMock.Expect(ctx, req).Return(res, nil)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				return nil
			},
			txManagerMockFunc: func(mc minimock.MockController) *dbMock.TxManagerMock {
				return dbMock.NewTxManagerMock(mc)
			},
		},
		{
			name: "User repo failure",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.GetListMock.Expect(ctx, req).Return(nil, repoErr)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				return nil
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
			userCacheMock := tt.userCacheMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)

			confMock := configMock.NewUserServiceConfigMock(mc)
			confMock.NoCacheMock.Optional().Set(func() bool { return true })

			service := userService.NewUserService(userRepoMock, userCacheMock, txManagerMock, confMock)

			servRes, err := service.GetList(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, servRes)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
