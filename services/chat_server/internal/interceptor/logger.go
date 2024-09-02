package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
)

// ClientLoggerInterceptor logs request by rpc client
func ClientLoggerInterceptor(ctx context.Context, method string, req, res any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	now := time.Now()
	err := invoker(ctx, method, req, res, cc, opts...)
	if err != nil {
		logger.Error("failed to request rpc client",
			zap.String("method", method),
			zap.String("client", cc.Target()),
			zap.Any("req", req),
			zap.Error(err),
		)
	}

	logger.Debug("rpc client request",
		zap.String("method", method),
		zap.String("client", cc.Target()),
		zap.Any("req", req),
		zap.Any("res", res),
		zap.Duration("duration", time.Since(now)),
	)

	return err
}
