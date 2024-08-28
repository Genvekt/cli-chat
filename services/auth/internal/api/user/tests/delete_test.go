package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	userImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/user"
	serviceMock "github.com/Genvekt/cli-chat/services/auth/internal/service/mocks"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *userApi.DeleteRequest
	}

	type userServiceMockFunc func(mc minimock.MockController) *serviceMock.UserServiceMock

	var (
		cxt        = context.Background()
		mc         = minimock.NewController(t)
		id         = gofakeit.Int64()
		serviceErr = fmt.Errorf("service error")
		req        = &userApi.DeleteRequest{Id: id}
		res        = &emptypb.Empty{}
	)

	tests := []struct {
		name            string
		args            args
		wand            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: cxt,
				req: req,
			},
			wand: res,
			err:  nil,
			userServiceMock: func(mc minimock.MockController) *serviceMock.UserServiceMock {
				mock := serviceMock.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(cxt, id).Return(nil)
				return mock
			},
		},
		{
			name: "Failure",
			args: args{
				ctx: cxt,
				req: req,
			},
			wand: nil,
			err:  serviceErr,
			userServiceMock: func(mc minimock.MockController) *serviceMock.UserServiceMock {
				mock := serviceMock.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(cxt, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userImpl.NewService(userServiceMock, nil)

			resp, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wand, resp)
			require.Equal(t, tt.err, err)
		})
	}
}
