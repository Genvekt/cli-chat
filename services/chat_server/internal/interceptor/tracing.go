package interceptor

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const traceIdHeader = "x-trace-id"

func ServerTracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if ok {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(traceIdHeader, spanContext.TraceID().String()))
		header := metadata.New(map[string]string{traceIdHeader: spanContext.TraceID().String()})
		err := grpc.SendHeader(ctx, header)
		if err != nil {
			return nil, err
		}
	}

	res, err := handler(ctx, req)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("err", err.Error())
	}

	return res, err
}

func ClientTracingInterceptor(ctx context.Context, method string, req, res any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, method)
	defer span.Finish()

	spanContext, ok := span.Context().(jaeger.SpanContext)
	if ok {
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(traceIdHeader, spanContext.TraceID().String()))
		//header := metadata.New(map[string]string{traceIdHeader: spanContext.TraceID().String()})
		//err := grpc.SendHeader(ctx, header)
		//if err != nil {
		//	return err
		//}
	}

	err := invoker(ctx, method, req, res, cc, opts...)
	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag("err", err.Error())
	}

	return err
}
