package token

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils"

	"github.com/golang-jwt/jwt"
)

var _ utils.TokenProvider = (*tokenProvider)(nil)

type tokenProvider struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenProvider creates token utils based on provided configuration
func NewTokenProvider(conf config.TokenProviderConfig) *tokenProvider {
	return &tokenProvider{
		secret: conf.Secret(),
		ttl:    conf.TTL(),
	}
}

// Generate creates new token for user
func (t *tokenProvider) Generate(ctx context.Context, user *model.User) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "generate token")
	defer span.Finish()

	claims := &model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t.ttl).Unix(),
		},
		ID:       user.ID,
		Username: user.Name,
		Role:     user.Role,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString(t.secret)
}

// Verify checks tha token was generated by this instance
func (t *tokenProvider) Verify(ctx context.Context, token string) (*model.UserClaims, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "verify token")
	defer span.Finish()

	jwtToken, err := jwt.ParseWithClaims(token, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected token signing method")
		}

		return t.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := jwtToken.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
