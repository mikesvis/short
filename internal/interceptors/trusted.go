package interceptors

import (
	"context"
	"slices"

	"github.com/mikesvis/short/internal/subnet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Проверяет доступ по маске подсети из заголовка X-Real-IP
func TrustedSubnetInterceptor(ts string) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		endPoints := []string{"/short.ShortService/GetStats"}
		if !slices.Contains(endPoints, info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		// Metadata пустая, но и настройка открыта для всех - можно
		if !ok || ts == "" {
			return handler(ctx, req)
		}

		// Metadata пустая, но настройка закрыта по маске - нельзя
		if !ok && ts != "" {
			return nil, status.Errorf(codes.PermissionDenied, "Forbidden area")
		}

		// Ключ не найден в Metadata - нельзя
		xRealIPSlice, ok := md["x-real-ip"]
		if !ok || len(xRealIPSlice) == 0 {
			return nil, status.Error(codes.PermissionDenied, "x-real-ip is not provided")
		}

		clientIP := xRealIPSlice[0]
		if !subnet.ValidateSubnet(clientIP, ts) {
			return nil, status.Error(codes.PermissionDenied, "x-real-ip is empty")
		}

		return handler(ctx, req)
	}
}
