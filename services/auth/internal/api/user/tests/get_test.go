package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	userImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/user"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	serviceMock "github.com/Genvekt/cli-chat/services/auth/internal/service/mocks"
)

func TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
		req *userApi.GetRequest
	}

	type userServiceMockFunc func(mc minimock.MockController) *serviceMock.UserServiceMock

	var (
		cxt       = context.Background()
		mc        = minimock.NewController(t)
		id        = gofakeit.Int64()
		name      = gofakeit.Username()
		email     = gofakeit.Email()
		roleUser  = userApi.UserRole_USER
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		userModel = &model.User{
			ID:        id,
			Name:      name,
			Email:     email,
			Role:      int(roleUser),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		userInfo = &userApi.UserInfo{
			Name:  name,
			Email: email,
			Role:  roleUser,
		}

		user = &userApi.User{
			Id:        id,
			Info:      userInfo,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}

		serviceErr = fmt.Errorf("service error")
		req        = &userApi.GetRequest{Id: id}
		res        = &userApi.GetResponse{User: user}
	)

	tests := []struct {
		name            string
		args            args
		wand            *userApi.GetResponse
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
				mock.GetMock.Expect(cxt, id).Return(userModel, nil)
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
				mock.GetMock.Expect(cxt, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := tt.userServiceMock(mc)
			api := userImpl.NewService(userServiceMock)

			resp, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wand, resp)
			require.Equal(t, tt.err, err)
		})
	}
}
