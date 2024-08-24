package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

const (
	noCacheEnv     = "USER_SERVICE_NO_CACHE"
	cacheTTLMinEnv = "USER_SERVICE_CACHE_TTL_MIN"

	defaultCacheTTL = time.Minute
)

var _ config.UserServiceConfig = (*userServiceConfigEnv)(nil)

type userServiceConfigEnv struct {
	noCache  bool
	cacheTTL time.Duration
}

// NewUserServiceConfigEnv reads all user service configurations from env
func NewUserServiceConfigEnv() (*userServiceConfigEnv, error) {
	ttl := defaultCacheTTL

	ttlEnv := os.Getenv(cacheTTLMinEnv)
	if ttlEnv != "" {
		ttlMin, err := strconv.Atoi(ttlEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid %s value: expected int, got %s", cacheTTLMinEnv, ttlEnv)
		}
		ttl = time.Duration(ttlMin) * time.Minute
	}

	noCache := false

	noCacheUsedEnv := os.Getenv(noCacheEnv)
	if noCacheUsedEnv != "" {
		noCacheUsed, err := strconv.ParseBool(noCacheUsedEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid %s value: expected bool, got %s", noCacheEnv, noCacheUsedEnv)
		}
		noCache = noCacheUsed
	}

	return &userServiceConfigEnv{
		cacheTTL: ttl,
		noCache:  noCache,
	}, nil
}

// CacheTTL returns the ttl of cache
func (s *userServiceConfigEnv) CacheTTL() time.Duration {
	return s.cacheTTL
}

// NoCache flag indicates that cache is not used
func (s *userServiceConfigEnv) NoCache() bool {
	return s.noCache
}

// UseCache flag indicates that cache is used
func (s *userServiceConfigEnv) UseCache() bool {
	return !s.noCache
}
