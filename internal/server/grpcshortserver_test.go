package server

import (
	_context "context"
	goerrors "errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/drivers/inmemory"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/interceptors"
	"github.com/mikesvis/short/internal/jwt"
	"github.com/mikesvis/short/internal/logger"
	pb "github.com/mikesvis/short/internal/proto"
	"github.com/mikesvis/short/internal/storage"
	mock_storage "github.com/mikesvis/short/mocks/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

func initTestServer(logger *zap.SugaredLogger, config *config.Config, service *ShortGRPCService) (*bufconn.Listener, *grpc.Server) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer(interceptors.GetInterceptors(logger, config))
	pb.RegisterShortServiceServer(server, service)
	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Fatalf("Server exited with error: %v", err)
		}
	}()

	return listener, server
}

func initBalalaykaTestStandByConfig(t *testing.T, config *config.Config) (*bufconn.Listener, *grpc.Server, *mock_storage.MockStorageDeleter) {
	// logger
	logger, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("Unable to init logger: %v", err)
	}
	// storage mocked
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storage := mock_storage.NewMockStorageDeleter(ctrl)
	// service
	service := NewShortService(storage, config)
	// server
	listener, server := initTestServer(logger, config, service)

	return listener, server, storage
}

func TestNewShortServer(t *testing.T) {
	logger, _ := logger.NewLogger()
	type args struct {
		storage storage.Storage
		config  *config.Config
	}
	tests := []struct {
		name string
		args args
		want *ShortGRPCService
	}{
		{
			name: "New server",
			args: args{
				inmemory.NewInMemory(logger),
				&config.Config{
					ServerAddress:     "localhost:8080",
					BaseURL:           "http://localhost:8080",
					FileStoragePath:   "",
					DatabaseDSN:       "",
					EnableHTTPS:       false,
					ServerKeyPath:     "",
					ServerCertPath:    "",
					GRPCServerAddress: "localhost:8082",
				},
			},
			want: &ShortGRPCService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShortService(tt.args.storage, tt.args.config)
			require.NotEmpty(t, got)
		})
	}
}

func TestShortGRPCService_SaveURL(t *testing.T) {
	config := &config.Config{
		GRPCServerAddress: ":5051",
		BaseURL:           "http://localhost",
	}
	listener, server, mockedStorage := initBalalaykaTestStandByConfig(t, config)

	// mocks start
	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("short")
	mockedStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(domain.URL{
		UserID: "DoomGuy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	}, nil)

	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("short")
	mockedStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(domain.URL{
		UserID: "Heretic",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	}, errors.ErrConflict)

	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("short")
	mockedStorage.EXPECT().Store(gomock.Any(), gomock.Any()).Return(domain.URL{
		UserID: "Unlucky guy",
		Full:   "http://www.yandex.ru/verylongpath",
		Short:  "short",
	}, goerrors.New("Internal server error"))
	// mocks end

	type args struct {
		ctx _context.Context
		in  *pb.SaveURLRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.SaveURLResponse
		wantErr bool
	}{
		{
			name: "Invalid URL",
			args: args{
				ctx: _context.Background(),
				in:  &pb.SaveURLRequest{Url: "http"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success save URL",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				in:  &pb.SaveURLRequest{Url: "http://www.yandex.ru/verylongpath"},
			},
			want: &pb.SaveURLResponse{
				ShortURL: config.BaseURL + `/short`,
			},
			wantErr: false,
		},
		{
			name: "Conflict save URL",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Heretic"),
				in:  &pb.SaveURLRequest{Url: "http://www.yandex.ru/verylongpath"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Save failed due to server",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Unlucky guy"),
				in:  &pb.SaveURLRequest{Url: "http://www.yandex.ru/verylongpath"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(_context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := pb.NewShortServiceClient(conn)
			response, err := client.SaveURL(tt.args.ctx, tt.args.in)

			if tt.wantErr {
				require.Error(t, err)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.ShortURL, response.ShortURL)
			}

		})
	}
	t.Cleanup(func() {
		server.Stop()
	})
}

func TestShortGRPCService_ShortenBatchURL(t *testing.T) {
	config := &config.Config{
		GRPCServerAddress: ":5051",
		BaseURL:           "http://localhost",
	}
	listener, server, mockedStorage := initBalalaykaTestStandByConfig(t, config)

	// mocks start
	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("shor1")
	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("shor2")
	mockedStorage.EXPECT().StoreBatch(gomock.Any(), gomock.Any()).Return(map[string]domain.URL{
		"1": {
			UserID: "DoomGuy",
			Full:   "http://www.yandex.ru/verylongpath1",
			Short:  "shor1",
		},
		"2": {
			UserID: "DoomGuy",
			Full:   "http://www.yandex.ru/verylongpath2",
			Short:  "shor2",
		},
	}, nil)

	mockedStorage.EXPECT().GetRandkey(uint(5)).Return("shor1")
	mockedStorage.EXPECT().StoreBatch(gomock.Any(), gomock.Any()).Return(map[string]domain.URL{}, goerrors.New("Internal server error"))
	// mocks end

	type args struct {
		ctx _context.Context
		in  *pb.ShortenBatchURLRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.ShortenBatchURLResponse
		wantErr bool
	}{
		{
			name: "Success save URLs",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "DoomGuy"),
				in: &pb.ShortenBatchURLRequest{
					Urls: []*pb.ShortenBatchURLRequestElement{
						{
							CorrelationId: "1",
							OriginalURL:   "http://www.yandex.ru/verylongpath1",
						},
						{
							CorrelationId: "2",
							OriginalURL:   "http://www.yandex.ru/verylongpath2",
						},
					},
				},
			},
			want: &pb.ShortenBatchURLResponse{
				Urls: []*pb.ShortenBatchURLResponseElement{
					{
						CorrelationId: "1",
						ShortURL:      config.BaseURL + `/shor1`,
					},
					{
						CorrelationId: "2",
						ShortURL:      config.BaseURL + `/shor2`,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Internal server error",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Unlucky guy"),
				in: &pb.ShortenBatchURLRequest{
					Urls: []*pb.ShortenBatchURLRequestElement{
						{
							CorrelationId: "1",
							OriginalURL:   "http://www.yandex.ru/verylongpath1",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(_context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := pb.NewShortServiceClient(conn)
			response, err := client.ShortenBatchURL(tt.args.ctx, tt.args.in)

			if tt.wantErr {
				require.Error(t, err)
			}

			if tt.want != nil {
				assert.ElementsMatch(t, tt.want.Urls, response.Urls)
			}

		})
	}
	t.Cleanup(func() {
		server.Stop()
	})
}

func TestShortGRPCService_GetURLByID(t *testing.T) {
	config := &config.Config{
		GRPCServerAddress: ":5051",
		BaseURL:           "http://localhost",
	}
	listener, server, mockedStorage := initBalalaykaTestStandByConfig(t, config)

	// mocks start
	mockedStorage.EXPECT().GetByShort(gomock.Any(), gomock.Any()).Return(domain.URL{}, goerrors.New("Internal server error"))
	mockedStorage.EXPECT().GetByShort(gomock.Any(), gomock.Any()).Return(domain.URL{}, nil)
	mockedStorage.EXPECT().GetByShort(gomock.Any(), gomock.Any()).Return(domain.URL{
		UserID:  "Doomguy",
		Full:    "http://ya.ru",
		Short:   "short1",
		Deleted: true,
	}, nil)
	mockedStorage.EXPECT().GetByShort(gomock.Any(), gomock.Any()).Return(domain.URL{
		UserID:  "Doomguy",
		Full:    "http://ya.ru",
		Short:   "short1",
		Deleted: false,
	}, nil)
	// mocks end

	type args struct {
		ctx _context.Context
		in  *pb.GetURLByIDRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetURLByIDResponse
		wantErr bool
	}{
		{
			name: "Getting failed due to server",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Unlucky guy"),
				in: &pb.GetURLByIDRequest{
					ShortURL: "this is not working",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Is deleted",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Doomguy"),
				in: &pb.GetURLByIDRequest{
					ShortURL: "short1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Is deleted",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Doomguy"),
				in: &pb.GetURLByIDRequest{
					ShortURL: "short1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success get by ID",
			args: args{
				ctx: _context.WithValue(_context.Background(), context.UserIDContextKey, "Doomguy"),
				in: &pb.GetURLByIDRequest{
					ShortURL: "short1",
				},
			},
			want: &pb.GetURLByIDResponse{
				OriginalURL: `http://ya.ru`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(_context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := pb.NewShortServiceClient(conn)
			response, err := client.GetURLByID(tt.args.ctx, tt.args.in)

			if tt.wantErr {
				require.Error(t, err)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.OriginalURL, response.OriginalURL)
			}

		})
	}
	t.Cleanup(func() {
		server.Stop()
	})
}

func TestShortGRPCService_GetURLByUser(t *testing.T) {
	config := &config.Config{
		GRPCServerAddress: ":5051",
		BaseURL:           "http://localhost",
	}
	listener, server, mockedStorage := initBalalaykaTestStandByConfig(t, config)
	expirationTime := time.Now().Add(jwt.TokenDuration)
	token, err := jwt.CreateTokenString("Doomguy", expirationTime)
	if err != nil {
		log.Fatalf("Failed create token %v", err)
	}

	// mocks start
	mockedStorage.EXPECT().GetUserURLs(gomock.Any(), gomock.Any()).Return([]domain.URL{}, goerrors.New("Internal server error"))
	mockedStorage.EXPECT().GetUserURLs(gomock.Any(), gomock.Any()).Return([]domain.URL{}, nil)
	mockedStorage.EXPECT().GetUserURLs(gomock.Any(), "Doomguy").Return([]domain.URL{
		{
			UserID:  "Doomguy",
			Short:   "short",
			Full:    "http://ya.ru",
			Deleted: false,
		},
	}, nil)
	// mocks end

	type args struct {
		ctx _context.Context
		in  *emptypb.Empty
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetURLByUserResponse
		wantErr bool
	}{
		{
			name: "Getting failed due to server",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{jwt.AuthorizationMDKeyName: token}),
				),
				in: &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Getting empty result",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{jwt.AuthorizationMDKeyName: token}),
				),
				in: &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success getting user URLs",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{jwt.AuthorizationMDKeyName: token}),
				),
				in: &emptypb.Empty{},
			},
			want: &pb.GetURLByUserResponse{
				Urls: []*pb.URLByUserResponseElement{
					{
						ShortURL:    config.BaseURL + `/short`,
						OriginalURL: `http://ya.ru`,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(_context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := pb.NewShortServiceClient(conn)
			response, err := client.GetURLByUser(tt.args.ctx, tt.args.in)

			if tt.wantErr {
				require.Error(t, err)
			}

			if tt.want != nil {
				assert.ElementsMatch(t, tt.want.Urls, response.Urls)
			}
		})
	}
	t.Cleanup(func() {
		server.Stop()
	})
}

func TestShortGRPCService_DeleteBatchURL(t *testing.T) {
	config := &config.Config{
		GRPCServerAddress: ":5051",
		BaseURL:           "http://localhost",
	}
	listener, server, mockedStorage := initBalalaykaTestStandByConfig(t, config)
	expirationTime := time.Now().Add(jwt.TokenDuration)
	token, err := jwt.CreateTokenString("Doomguy", expirationTime)
	if err != nil {
		log.Fatalf("Failed create token %v", err)
	}

	// mocks start
	mockedStorage.EXPECT().DeleteBatch(gomock.Any(), "Doomguy", gomock.Any()).Return()
	// mocks end

	type args struct {
		ctx _context.Context
		in  *pb.DeleteBatchURLRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty input",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{jwt.AuthorizationMDKeyName: token}),
				),
				in: &pb.DeleteBatchURLRequest{},
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{jwt.AuthorizationMDKeyName: token}),
				),
				in: &pb.DeleteBatchURLRequest{
					Urls: []string{"some", "thing"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(_context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := pb.NewShortServiceClient(conn)
			_, err = client.DeleteBatchURL(tt.args.ctx, tt.args.in)

			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
	t.Cleanup(func() {
		server.Stop()
	})
}

func TestShortGRPCService_GetStats(t *testing.T) {
	config := &config.Config{
		GRPCServerAddress: ":5051",
		BaseURL:           "http://localhost",
		TrustedSubnet:     "127.0.0.1/24",
	}
	listener, server, mockedStorage := initBalalaykaTestStandByConfig(t, config)
	expirationTime := time.Now().Add(jwt.TokenDuration)
	token, err := jwt.CreateTokenString("Doomguy", expirationTime)
	if err != nil {
		log.Fatalf("Failed create token %v", err)
	}

	// mocks start
	mockedStorage.EXPECT().GetStats(gomock.Any()).Return(domain.Stats{}, goerrors.New("Failed"))
	mockedStorage.EXPECT().GetStats(gomock.Any()).Return(domain.Stats{
		Users: 1,
		URLs:  1,
	}, nil)
	// mocks end

	type args struct {
		ctx _context.Context
		in  *emptypb.Empty
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.GetStatsResponse
		wantErr bool
	}{
		{
			name: "Closed",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{
						"x-real-ip":                "143.0.0.4",
						jwt.AuthorizationMDKeyName: token,
					}),
				),
				in: &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Internal server error",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{
						"x-real-ip":                "127.0.0.4",
						jwt.AuthorizationMDKeyName: token,
					}),
				),
				in: &emptypb.Empty{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				ctx: metadata.NewOutgoingContext(
					_context.Background(),
					metadata.New(map[string]string{
						"x-real-ip":                "127.0.0.4",
						jwt.AuthorizationMDKeyName: token,
					}),
				),
				in: &emptypb.Empty{},
			},
			want: &pb.GetStatsResponse{
				Urls:  1,
				Users: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(func(_context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			client := pb.NewShortServiceClient(conn)
			response, err := client.GetStats(tt.args.ctx, tt.args.in)

			if tt.wantErr {
				require.Error(t, err)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.Urls, response.Urls)
				assert.Equal(t, tt.want.Users, response.Users)
			}
		})
	}
	t.Cleanup(func() {
		server.Stop()
	})
}
