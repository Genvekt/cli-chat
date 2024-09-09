package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/Genvekt/cli-chat/services/auth/internal/model"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
)

// Check validates access permissions for endpoint
func (s *accessService) Check(ctx context.Context, accessToken string, endpoint string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "check access")
	defer span.Finish()

	claims, err := s.accessTokenProvider.Verify(ctx, accessToken)
	if err != nil {
		return false, fmt.Errorf("access token is invalid")
	}

	return s.hasAccess(ctx, claims, endpoint)
}

func (s *accessService) hasAccess(ctx context.Context, claims *model.UserClaims, endpoint string) (bool, error) {
	rule, err := s.accessRepo.GetEndpointAccessRule(ctx, endpoint)
	if err != nil {
		if errors.Is(err, repository.ErrRuleNotFound) {
			// No rules specified for endpoint allows access for all
			return true, nil
		}
		return false, err
	}

	if rule.HasRole(claims.Role) {
		return true, nil
	}

	return false, nil
}
