package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

const (
	defaultCacheTTL = time.Minute
)

var _ config.UserServiceConfig = (*userServiceConfigEnv)(nil)

type userServiceConfigEnv struct {
	cacheTTL time.Duration
}

// NewUserServiceConfigEnv reads all user service configurations from env
func NewUserServiceConfigEnv() (*userServiceConfigEnv, error) {
	ttl := defaultCacheTTL

	ttlEnv := os.Getenv("USER_SERVICE_CACHE_TTL_MIN")
	if ttlEnv != "" {
		ttlMin, err := strconv.Atoi(ttlEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid USER_SERVICE_CACHE_TTL_MIN value: expected int, got %s", ttlEnv)
		}
		ttl = time.Duration(ttlMin) * time.Minute
	}
	return &userServiceConfigEnv{
		cacheTTL: ttl,
	}, nil
}

func (s *userServiceConfigEnv) CacheTTL() time.Duration {
	return s.cacheTTL
}
