// Пакет приложения сокращателя ссылок.
package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/interceptors"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/middleware"
	pb "github.com/mikesvis/short/internal/proto"
	"github.com/mikesvis/short/internal/server"
	"github.com/mikesvis/short/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
)

// App - стуктура приложения с конфигом, логгером, storage и роутером.
type App struct {
	config     *config.Config
	logger     *zap.SugaredLogger
	storage    storage.Storage
	router     *chi.Mux
	httpServer *http.Server
	grpcServer *grpc.Server
}

// Конструктор приложения, здесь инициализируются все зависимости:
// конфиг приложения, логгер, storage, роутер. Также здесь регистрируются middleware приложения.
func New(config *config.Config) *App {
	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	storage, err := storage.NewStorage(config, logger)
	if err != nil {
		panic(err)
	}

	// инфраструктура для HTTP
	handler := server.NewHandler(config, storage)
	router := server.NewRouter(
		handler,
		middleware.RequestResponseLogger(logger),
		middleware.GZip(
			[]string{
				"application/json",
				"text/html",
				"application/x-gzip",
			},
		),
	)

	httpServer := &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}

	// инфраструктура для GRPC
	interceptors := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			interceptors.RequestResponseLoggerInterceptor(logger),
			interceptors.SignInInterceptor,
			interceptors.AuthInterceptor,
			interceptors.TrustedSubnetInterceptor(config.TrustedSubnet),
		),
	)

	var grpcServer *grpc.Server
	if config.EnableHTTPS {
		creds, err := credentials.NewServerTLSFromFile(config.ServerCertPath, config.ServerKeyPath)
		if err != nil {
			logger.Fatalf("Failed to init creds for gRPC: %v", err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds), interceptors)
	} else {
		grpcServer = grpc.NewServer(interceptors)
	}

	shortGRPCServer := server.NewShortServer(storage, config)
	pb.RegisterShortServiceServer(grpcServer, shortGRPCServer)

	return &App{
		config,
		logger,
		storage,
		router,
		httpServer,
		grpcServer,
	}
}

// Запуск приложения.
func (a *App) Run() error {
	a.logger.Infow("Config initialized", "config", a.config)
	if _, isCloser := a.storage.(storage.StorageCloser); isCloser {
		defer a.storage.(storage.StorageCloser).Close()
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Запускаем GRPC сервер
	go func() {
		listener, err := net.Listen("tcp", a.config.GRPCServerAddress)
		if err != nil {
			a.logger.Fatalf("Failed to listen for gRPC: %v", err)
		}
		if err := a.grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Запускаем HTTP сервер
	go func() {
		if a.config.EnableHTTPS {
			if err := a.httpServer.ListenAndServeTLS(a.config.ServerCertPath, a.config.ServerKeyPath); err != http.ErrServerClosed {
				a.logger.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			if err := a.httpServer.ListenAndServe(); err != http.ErrServerClosed {
				a.logger.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
		return err
	}
	log.Println("Server exited properly")
	return nil
}
