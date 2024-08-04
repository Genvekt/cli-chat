package redis

import (
	"context"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"

	modelRepo "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/redis/model"

	"github.com/Genvekt/cli-chat/libraries/cache_client/pkg/cache"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository/user/redis/converter"
)

var _ repository.UserCache = (*userCacheRedis)(nil)

type userCacheRedis struct {
	client cache.RedisClient
}

// NewUserCacheRedis obtains instance of user redis cache
func NewUserCacheRedis(client cache.RedisClient) *userCacheRedis {
	return &userCacheRedis{
		client: client,
	}
}

func (r *userCacheRedis) Get(ctx context.Context, id int64) (*model.User, error) {
	idStr := strconv.FormatInt(id, 10)
	values, err := r.client.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, repository.ErrUserNotFound
	}

	var user modelRepo.User
	err = redigo.ScanStruct(values, &user)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (r *userCacheRedis) Set(ctx context.Context, user *model.User) error {
	userRepo := converter.ToRepoFromUser(user)

	idStr := strconv.FormatInt(userRepo.ID, 10)
	err := r.client.HashSet(ctx, idStr, userRepo)
	if err != nil {
		return err
	}

	return nil
}

func (r *userCacheRedis) Expire(ctx context.Context, id int64, timeout time.Duration) error {
	idStr := strconv.FormatInt(id, 10)

	err := r.client.Expire(ctx, idStr, timeout)
	if err != nil {
		return err
	}

	return nil
}
