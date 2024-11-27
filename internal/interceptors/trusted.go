package interceptors

import (
	"context"
	"net"
	"slices"

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
		// X-Real-IP пустой - нелья
		if clientIP == "" {
			return nil, status.Error(codes.PermissionDenied, "x-real-ip is empty")
		}

		// Ошибка при парсинге - нельзя
		_, cidr, err := net.ParseCIDR(ts)
		if err != nil {
			return nil, status.Error(codes.Internal, "Subnet mask parsing error")
		}

		// IP не входит в доверенную сеть - нельзя
		if !cidr.Contains(net.ParseIP(clientIP)) {
			return nil, status.Error(codes.PermissionDenied, "Forbidden area")
		}

		// МОЖНА!
		return handler(ctx, req)
	}
}
