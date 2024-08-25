package access

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	accessApi "github.com/Genvekt/cli-chat/libraries/api/access/v1"
)

const (
	authorisationHeader = "authorization"
	authPrefix          = "Bearer "
)

// Check validates client has access to endpoint
func (s *Service) Check(ctx context.Context, req *accessApi.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not provided")
	}
	authHeader := md.Get(authorisationHeader)
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("%s header not provided", authorisationHeader)
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, fmt.Errorf("invalid %s header format", authorisationHeader)
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	hasAccess, err := s.accessService.Check(ctx, accessToken, req.GetEndpointAddress())
	if err != nil || !hasAccess {
		return nil, fmt.Errorf("access denied")
	}

	return &emptypb.Empty{}, nil
}
