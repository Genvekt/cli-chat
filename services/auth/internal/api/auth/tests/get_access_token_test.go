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

func TestGetAccessToken(t *testing.T) {
	type args struct {
		ctx context.Context
		req *authApi.GetAccessTokenRequest
	}

	type authServiceMockFunc func(mc minimock.MockController) *serviceMock.AuthServiceMock

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		oldToken = gofakeit.HackerPhrase()
		newToken = gofakeit.HackerPhrase()

		serviceErr = fmt.Errorf("service error")

		req = &authApi.GetAccessTokenRequest{
			RefreshToken: oldToken,
		}

		res = &authApi.GetAccessTokenResponse{
			AccessToken: newToken,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                *authApi.GetAccessTokenResponse
		err                 error
		authServiceMockFunc authServiceMockFunc
	}{
		{
			name: "Refresh update success",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMockFunc: func(mc minimock.MockController) *serviceMock.AuthServiceMock {
				mock := serviceMock.NewAuthServiceMock(mc)
				mock.GetAccessTokenMock.Expect(ctx, oldToken).Return(newToken, nil)
				return mock
			},
		},
		{
			name: "Refresh update failed",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			authServiceMockFunc: func(mc minimock.MockController) *serviceMock.AuthServiceMock {
				mock := serviceMock.NewAuthServiceMock(mc)
				mock.GetAccessTokenMock.Expect(ctx, oldToken).Return("", serviceErr)
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

			apiRes, err := api.GetAccessToken(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.want, apiRes)
			require.Equal(t, tt.err, err)
		})
	}
}
