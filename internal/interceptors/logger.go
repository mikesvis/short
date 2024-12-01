package interceptors

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Логирование запросов и времени обработки.
func RequestResponseLoggerInterceptor(log *zap.SugaredLogger) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		method := info.FullMethod

		// не уверен насколько это правильно
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		log.Infow(
			"Incoming GRPC request",
			"method", method,
			"status", status.Code(err),
			"duration", duration,
		)

		return resp, err
	}
}
