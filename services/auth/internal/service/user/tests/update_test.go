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

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	repoMock "github.com/Genvekt/cli-chat/services/auth/internal/repository/mocks"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

func TestUpdateWithCache(t *testing.T) {
	type args struct {
		ctx context.Context
		dto *model.UserUpdateDTO
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type userCacheMockFunc func(mc minimock.MockController) *repoMock.UserCacheMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock
	type configMockFunc func(mc minimock.MockController) *configMock.UserServiceConfigMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		updatedName  = gofakeit.Name()
		updatedEmail = gofakeit.Email()
		updatedRole  = 1
		createdAt    = gofakeit.Date()

		userBeforeUpdate = &model.User{
			ID:        id,
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      2,
			CreatedAt: createdAt,
		}
		updateDto = &model.UserUpdateDTO{
			ID:    id,
			Name:  &updatedName,
			Email: &updatedEmail,
			Role:  &updatedRole,
		}
		repoErr = fmt.Errorf("repo error")
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
			name: "Success",
			args: args{
				ctx: ctx,
				dto: updateDto,
			},
			err: nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Set(
					func(_ context.Context, userID int64, updateFunc func(user *model.User) error) (err error) {

						require.Equal(t, userID, id)
						err = updateFunc(userBeforeUpdate)
						if err != nil {
							return err
						}

						// Check that update was correct
						require.Equal(t, userBeforeUpdate.ID, id)
						require.Equal(t, userBeforeUpdate.Name, updatedName)
						require.Equal(t, userBeforeUpdate.Email, updatedEmail)
						require.Equal(t, userBeforeUpdate.Role, updatedRole)
						require.Equal(t, userBeforeUpdate.CreatedAt, createdAt)

						return nil
					},
				)
				return mock
			},
			userCacheMockFunc: func(mc minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(false)
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
				dto: updateDto,
			},
			err: repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Return(repoErr)
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
				dto: updateDto,
			},
			err: repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Return(nil)
				return mock
			},
			userCacheMockFunc: func(mc minimock.MockController) *repoMock.UserCacheMock {
				mock := repoMock.NewUserCacheMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
				return mock
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(false)
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

			err := service.Update(tt.args.ctx, tt.args.dto)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}

func TestUpdateWithoutCache(t *testing.T) {
	type args struct {
		ctx context.Context
		dto *model.UserUpdateDTO
	}

	type userRepoMockFunc func(mc minimock.MockController) *repoMock.UserRepositoryMock
	type userCacheMockFunc func(mc minimock.MockController) *repoMock.UserCacheMock
	type txManagerMockFunc func(mc minimock.MockController) *dbMock.TxManagerMock
	type configMockFunc func(mc minimock.MockController) *configMock.UserServiceConfigMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id           = gofakeit.Int64()
		updatedName  = gofakeit.Name()
		updatedEmail = gofakeit.Email()
		updatedRole  = 1
		createdAt    = gofakeit.Date()

		userBeforeUpdate = &model.User{
			ID:        id,
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			Role:      2,
			CreatedAt: createdAt,
		}
		updateDto = &model.UserUpdateDTO{
			ID:    id,
			Name:  &updatedName,
			Email: &updatedEmail,
			Role:  &updatedRole,
		}
		repoErr = fmt.Errorf("repo error")
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
			name: "Success",
			args: args{
				ctx: ctx,
				dto: updateDto,
			},
			err: nil,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Set(
					func(_ context.Context, id int64, updateFunc func(user *model.User) error) (err error) {
						err = updateFunc(userBeforeUpdate)
						if err != nil {
							return err
						}

						// Check that update was correct
						require.Equal(t, userBeforeUpdate.ID, id)
						require.Equal(t, userBeforeUpdate.Name, updatedName)
						require.Equal(t, userBeforeUpdate.Email, updatedEmail)
						require.Equal(t, userBeforeUpdate.Role, updatedRole)
						require.Equal(t, userBeforeUpdate.CreatedAt, createdAt)

						return nil
					},
				)
				return mock
			},
			userCacheMockFunc: func(_ minimock.MockController) *repoMock.UserCacheMock {
				return nil
			},
			configMockFunc: func(mc minimock.MockController) *configMock.UserServiceConfigMock {
				mock := configMock.NewUserServiceConfigMock(mc)
				mock.NoCacheMock.Return(true)
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
				dto: updateDto,
			},
			err: repoErr,
			userRepoMockFunc: func(mc minimock.MockController) *repoMock.UserRepositoryMock {
				mock := repoMock.NewUserRepositoryMock(mc)
				mock.UpdateMock.Return(repoErr)
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := tt.userRepoMockFunc(mc)
			userCacheMock := tt.userCacheMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)
			confMock := tt.configMockFunc(mc)
			service := userService.NewUserService(userRepoMock, userCacheMock, txManagerMock, confMock)

			err := service.Update(tt.args.ctx, tt.args.dto)
			if tt.err == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
