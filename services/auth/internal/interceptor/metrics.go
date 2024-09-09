package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/Genvekt/cli-chat/services/auth/internal/metric"
)

// MetricsInterceptor records amount and duration of requests
func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	startTime := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(startTime)

	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTimeObserve("error", duration.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.HistogramResponseTimeObserve("success", duration.Seconds())
	}

	return res, err
}
