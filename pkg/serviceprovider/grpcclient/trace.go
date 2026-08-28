package grpcclient

import (
	"context"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/message"
	"time"

	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// getTraceId 获取追踪ID
func getTraceId(ctx context.Context) string {
	if id := ctxkey.GetValue(ctx, ctxkey.TraceIdKey); id != nil {
		if s, ok := id.(string); ok && s != "" {
			return s
		}
	}
	return "unknown"
}

// publishGrpcTrace 记录grpc调用调试信息
func publishGrpcTrace(ctx context.Context, method string, req, resp any, err error, ms float64) {
	message.NewEvent().Publish(debugger.TopicGrpc, debugger.GrpcEvent{
		TraceId:  getTraceId(ctx),
		Method:   method,
		Request:  req,
		Response: resp,
		Code:     status.Code(err).String(),
		Ms:       ms,
	})
}

// unaryClientInterceptor 客户端一元拦截器
func unaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpclib.ClientConn, invoker grpclib.UnaryInvoker, opts ...grpclib.CallOption) error {
	traceId := getTraceId(ctx)
	if traceId != "" && traceId != "unknown" {
		ctx = metadata.AppendToOutgoingContext(ctx, "trace-id", traceId)
	}

	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	publishGrpcTrace(ctx, method, req, reply, err, float64(time.Since(start).Nanoseconds())/1e6)
	return err
}

// unaryServerInterceptor 服务端一元拦截器
func unaryServerInterceptor(ctx context.Context, req any, info *grpclib.UnaryServerInfo, handler grpclib.UnaryHandler) (any, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("trace-id"); len(values) > 0 && values[0] != "" {
			ctx = ctxkey.WithValue(ctx, ctxkey.TraceIdKey, values[0])
		}
	}

	return handler(ctx, req)
}
