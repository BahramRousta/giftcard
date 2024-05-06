package utils

//import (
//	"context"
//	"go.opentelemetry.io/otel/trace"
//	"go.uber.org/zap"
//	"google.golang.org/grpc"
//	"giftCard/pkg/logger"
//)

//func GrpcLogInjector(log *zap.Logger) grpc.UnaryServerInterceptor {
//	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
//		traceId := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
//
//		lg := log.With(
//			zap.Any("TRACE.ID", traceId))
//
//		ctx = logger.ToContext(ctx, lg)
//
//		return handler(ctx, req)
//	}
//}
