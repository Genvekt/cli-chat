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

func TestGetList(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *userApi.GetListRequest
	}

	type userServiceMockFunc func(mc minimock.MockController) *serviceMock.UserServiceMock

	var (
		cxt = context.Background()
		mc  = minimock.NewController(t)

		userModel1 = &model.User{
			ID:        gofakeit.Int64(),
			Name:      gofakeit.Username(),
			Email:     gofakeit.Email(),
			Role:      int(userApi.UserRole_USER),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
		userModel2 = &model.User{
			ID:        gofakeit.Int64(),
			Name:      gofakeit.Username(),
			Email:     gofakeit.Email(),
			Role:      int(userApi.UserRole_USER),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		user1 = &userApi.User{
			Id: userModel1.ID,
			Info: &userApi.UserInfo{
				Name:  userModel1.Name,
				Email: userModel1.Email,
				Role:  userApi.UserRole(userModel1.Role),
			},
			CreatedAt: timestamppb.New(userModel1.CreatedAt),
			UpdatedAt: timestamppb.New(userModel1.UpdatedAt),
		}

		user2 = &userApi.User{
			Id: userModel2.ID,
			Info: &userApi.UserInfo{
				Name:  userModel2.Name,
				Email: userModel2.Email,
				Role:  userApi.UserRole(userModel2.Role),
			},
			CreatedAt: timestamppb.New(userModel2.CreatedAt),
			UpdatedAt: timestamppb.New(userModel2.UpdatedAt),
		}

		serviceErr = fmt.Errorf("service error")
		names      = []string{userModel1.Name, userModel2.Name, gofakeit.Username()}
		req        = &userApi.GetListRequest{Names: names}
		res        = &userApi.GetListResponse{Users: []*userApi.User{user1, user2}}
	)

	tests := []struct {
		name            string
		args            args
		wand            *userApi.GetListResponse
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
				mock.GetListMock.Expect(cxt, names).Return([]*model.User{userModel1, userModel2}, nil)
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
				mock.GetListMock.Expect(cxt, names).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userImpl.NewService(userServiceMock)

			resp, err := api.GetList(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.wand, resp)
			require.Equal(t, tt.err, err)
		})
	}
}
