package interceptors

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mikesvis/short/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func GetInterceptors(l *zap.SugaredLogger, c *config.Config) grpc.ServerOption {
	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			RequestResponseLoggerInterceptor(l),
			SignInInterceptor,
			AuthInterceptor,
			TrustedSubnetInterceptor(c.TrustedSubnet),
		),
	)
}
