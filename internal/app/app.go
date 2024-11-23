// Пакет приложения сокращателя ссылок.
package app

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/logger"
	"github.com/mikesvis/short/internal/middleware"
	"github.com/mikesvis/short/internal/server"
	"github.com/mikesvis/short/internal/storage"
	"go.uber.org/zap"
)

// App - стуктура приложения с конфигом, логгером, storage и роутером.
type App struct {
	config  *config.Config
	logger  *zap.SugaredLogger
	storage storage.Storage
	router  *chi.Mux
	server  *http.Server
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

	server := &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	}

	return &App{
		config,
		logger,
		storage,
		router,
		server,
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

	go func() {
		if a.config.EnableHTTPS {
			if err := a.server.ListenAndServeTLS(a.config.ServerCertPath, a.config.ServerKeyPath); err != http.ErrServerClosed {
				a.logger.Fatalf("Failed to start HTTPS server: %v", err)
			}
		} else {
			if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
				a.logger.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
		return err
	}
	log.Println("Server exited properly")
	return nil
}
