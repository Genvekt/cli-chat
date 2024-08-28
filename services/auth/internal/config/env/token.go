package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

const (
	refreshTokenTTLEnv    = "REFRESH_TOKEN_TTL_MIN" //nolint:gosec
	refreshTokenSecretEnv = "REFRESH_TOKEN_SECRET"  //nolint:gosec

	accessTokenTTLEnv    = "ACCESS_TOKEN_TTL_MIN" //nolint:gosec
	accessTokenSecretEnv = "ACCESS_TOKEN_SECRET"  //nolint:gosec
)

var _ config.TokenProviderConfig = (*tokenProviderConfig)(nil)

type tokenProviderConfig struct {
	ttl    time.Duration
	secret []byte
}

// NewRefreshTokenProviderConfig provides token configuration for refresh token
func NewRefreshTokenProviderConfig() (*tokenProviderConfig, error) {
	secret := os.Getenv(refreshTokenSecretEnv)
	if secret == "" {
		return nil, fmt.Errorf("environment variable %s not set", refreshTokenSecretEnv)
	}

	ttlStr := os.Getenv(refreshTokenTTLEnv)
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s not set, expecred number", refreshTokenTTLEnv)
	}

	return &tokenProviderConfig{
		ttl:    time.Duration(ttl) * time.Minute,
		secret: []byte(secret),
	}, nil
}

// NewAccessTokenProviderConfig provides token configuration for access token
func NewAccessTokenProviderConfig() (*tokenProviderConfig, error) {
	secret := os.Getenv(accessTokenSecretEnv)
	if secret == "" {
		return nil, fmt.Errorf("environment variable %s not set", accessTokenSecretEnv)
	}

	ttlStr := os.Getenv(accessTokenTTLEnv)
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s not set, expecred number", accessTokenTTLEnv)
	}

	return &tokenProviderConfig{
		ttl:    time.Duration(ttl) * time.Minute,
		secret: []byte(secret),
	}, nil
}

// TTL indicates the period when token is valid
func (t *tokenProviderConfig) TTL() time.Duration {
	return t.ttl
}

// Secret used to sign token and check its validity
func (t *tokenProviderConfig) Secret() []byte {
	return t.secret
}
