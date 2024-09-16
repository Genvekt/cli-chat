package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
)

// LogInterceptor logs processed gRPC requests
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()
	res, err := handler(ctx, req)
	if err != nil {
		logger.Error("failed to handle request",
			zap.String("method", info.FullMethod),
			zap.Any("req", req),
			zap.Error(err),
		)
	}

	logger.Debug("request",
		zap.String("method", info.FullMethod),
		zap.Any("req", req),
		zap.Any("res", res),
		zap.Duration("duration", time.Since(now)),
	)

	return res, err
}
