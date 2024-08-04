package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	userImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/user"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	serviceMock "github.com/Genvekt/cli-chat/services/auth/internal/service/mocks"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *userApi.UpdateRequest
	}

	type userServiceMockFunc func(mc minimock.MockController) *serviceMock.UserServiceMock

	var (
		cxt      = context.Background()
		mc       = minimock.NewController(t)
		id       = gofakeit.Int64()
		name     = gofakeit.Username()
		email    = gofakeit.Email()
		roleUser = userApi.UserRole_USER

		serviceErr = fmt.Errorf("service error")
		req        = &userApi.UpdateRequest{
			Id:    id,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
			Role:  &roleUser,
		}
		res = &emptypb.Empty{}
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
				mock.UpdateMock.Set(func(ctx context.Context, id int64, updateFunc func(user *model.User) error) (err error) {
					user := &model.User{ID: id}
					err = updateFunc(user)
					if err != nil {
						return err
					}
					require.Equal(t, user.ID, id)
					require.Equal(t, user.Name, name)
					require.Equal(t, user.Email, email)
					require.Equal(t, user.Role, int(roleUser))
					return nil
				})
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
				mock.UpdateMock.Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userImpl.NewService(userServiceMock)

			resp, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wand, resp)
			require.Equal(t, tt.err, err)
		})
	}
}
