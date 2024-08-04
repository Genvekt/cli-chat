package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	userImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/user"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	serviceMock "github.com/Genvekt/cli-chat/services/auth/internal/service/mocks"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *userApi.CreateRequest
	}

	type userServiceMockFunc func(mc minimock.MockController) *serviceMock.UserServiceMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Username()
		email     = gofakeit.Email()
		roleUser  = userApi.UserRole_USER
		password1 = gofakeit.Password(true, true, true, false, false, 15)
		password2 = gofakeit.Password(true, true, true, false, false, 15)

		serviceErr = fmt.Errorf("service error")

		userInfo = &userApi.UserInfo{
			Name:  name,
			Email: email,
			Role:  roleUser,
		}

		reqSamePass = &userApi.CreateRequest{
			Info:            userInfo,
			Password:        password1,
			PasswordConfirm: password1,
		}

		reqDiffPass = &userApi.CreateRequest{
			Info:            userInfo,
			Password:        password1,
			PasswordConfirm: password2,
		}

		userModel = &model.User{
			Name:  name,
			Email: email,
			Role:  int(roleUser),
		}

		res = &userApi.CreateResponse{
			Id: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *userApi.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Passwords match",
			args: args{
				ctx: ctx,
				req: reqSamePass,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc minimock.MockController) *serviceMock.UserServiceMock {
				mock := serviceMock.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userModel).Return(id, nil)
				return mock
			},
		},
		{
			name: "Passwords do not match",
			args: args{
				ctx: ctx,
				req: reqDiffPass,
			},
			want: nil,
			err:  userImpl.ErrPasswordsNotMatch,
			userServiceMock: func(mc minimock.MockController) *serviceMock.UserServiceMock {
				mock := serviceMock.NewUserServiceMock(mc)
				return mock
			},
		},
		{
			name: "Failure",
			args: args{
				ctx: ctx,
				req: reqSamePass,
			},
			want: nil,
			err:  serviceErr,
			userServiceMock: func(mc minimock.MockController) *serviceMock.UserServiceMock {
				mock := serviceMock.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userModel).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userImpl.NewService(userServiceMock)

			apiRes, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, err)
		})
	}
}
