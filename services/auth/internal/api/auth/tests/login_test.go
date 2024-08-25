package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
	authImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/auth"
	serviceMock "github.com/Genvekt/cli-chat/services/auth/internal/service/mocks"
)

func TestLogin(t *testing.T) {
	type args struct {
		ctx context.Context
		req *authApi.LoginRequest
	}

	type authServiceMockFunc func(mc minimock.MockController) *serviceMock.AuthServiceMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username = gofakeit.Username()
		password = gofakeit.Password(true, true, true, false, false, 15)
		token    = gofakeit.HackerPhrase()

		serviceErr = fmt.Errorf("service error")

		req = &authApi.LoginRequest{
			Username: username,
			Password: password,
		}

		res = &authApi.LoginResponse{
			RefreshToken: token,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *authApi.LoginResponse
		err                 error
		authServiceMockFunc authServiceMockFunc
	}{
		{
			name: "Login success",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMockFunc: func(mc minimock.MockController) *serviceMock.AuthServiceMock {
				mock := serviceMock.NewAuthServiceMock(mc)
				mock.LoginMock.Expect(ctx, username, password).Return(token, nil)
				return mock
			},
		},
		{
			name: "Login failure",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMockFunc: func(mc minimock.MockController) *serviceMock.AuthServiceMock {
				mock := serviceMock.NewAuthServiceMock(mc)
				mock.LoginMock.Expect(ctx, username, password).Return("", serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := tt.authServiceMockFunc(mc)
			api := authImpl.NewService(authServiceMock)

			apiRes, err := api.Login(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, err)
		})
	}
}
