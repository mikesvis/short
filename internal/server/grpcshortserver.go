package server

import (
	_context "context"
	_errors "errors"

	"github.com/mikesvis/short/internal/config"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/domain"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/keygen"
	pb "github.com/mikesvis/short/internal/proto"
	"github.com/mikesvis/short/internal/storage"
	"github.com/mikesvis/short/pkg/urlformat"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Сервер grpc для сервиса
type ShortGRPCService struct {
	pb.UnimplementedShortServiceServer
	storage storage.Storage
	config  *config.Config
}

// Инициализация сервера grpc для сервиса
func NewShortService(storage storage.Storage, config *config.Config) *ShortGRPCService {
	return &ShortGRPCService{storage: storage, config: config}
}

// Обработка /short.ShortService/SaveURL
// Запись сокращенного URL в условную "базу" если нет такого ключа
func (s *ShortGRPCService) SaveURL(ctx _context.Context, in *pb.SaveURLRequest) (*pb.SaveURLResponse, error) {
	var response pb.SaveURLResponse

	URL := string(in.Url)
	err := urlformat.ValidateURL(URL)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Invalid URL")
	}

	URL = urlformat.SanitizeURL(URL)

	item := domain.URL{
		UserID: ctx.Value(context.UserIDContextKey).(string),
		Full:   URL,
		Short:  s.storage.GetRandkey(keygen.KeyLength),
	}

	item, err = s.storage.Store(ctx, item)
	if err != nil && !_errors.Is(err, errors.ErrConflict) {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if err == nil {
		response.ShortURL = urlformat.FormatURL(string(s.config.BaseURL), item.Short)

		return &response, nil
	}

	return nil, status.Error(codes.Unknown, "Conflict")
}

// Обработка /short.ShortService/ShortenBatchURL
// Запись сокращенного URL в условную "базу" если нет такого ключа
func (s *ShortGRPCService) ShortenBatchURL(ctx _context.Context, in *pb.ShortenBatchURLRequest) (*pb.ShortenBatchURLResponse, error) {
	// генерим потенциальные domain.URL на сохранение с новым Short
	pack := make(map[string]domain.URL)
	for _, v := range in.Urls {
		pack[string(v.CorrelationId)] = domain.URL{
			UserID: ctx.Value(context.UserIDContextKey).(string),
			Full:   string(v.OriginalURL),
			Short:  s.storage.GetRandkey(keygen.KeyLength),
		}
	}

	// domain.URL.Short в процессе сохранения поменяем на старый если такой domain.URL.Full уже есть
	stored, err := s.storage.StoreBatch(ctx, pack)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	response := &pb.ShortenBatchURLResponse{}
	for k, v := range stored {
		responseElement := &pb.ShortenBatchURLResponseElement{
			CorrelationId: k,
			ShortURL:      urlformat.FormatURL(string(s.config.BaseURL), v.Short),
		}
		response.Urls = append(response.Urls, responseElement)
	}

	return response, nil
}

// Обработка /short.ShortService/GetURLByID
// GetURLByID получение полной ссылки по сокращенной
func (s *ShortGRPCService) GetURLByID(ctx _context.Context, in *pb.GetURLByIDRequest) (*pb.GetURLByIDResponse, error) {
	item, err := s.storage.GetByShort(ctx, in.ShortURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if (item == domain.URL{}) {
		return nil, status.Error(codes.Unknown, "Full URL is not found for provided short")
	}

	if item.Deleted {
		return nil, status.Error(codes.NotFound, "Resourse is gone")
	}

	return &pb.GetURLByIDResponse{OriginalURL: item.Full}, nil
}

// Обработка /short.ShortService/GetURLByUser
// GetURLByUser получение сокращенных ссылок пользователя
func (s *ShortGRPCService) GetURLByUser(ctx _context.Context, in *emptypb.Empty) (*pb.GetURLByUserResponse, error) {
	// тут умышленно ctx, ctx.Value
	// 1ый аргумент - контекст, 2ой аргумент - само значение ID пользователя
	userID := ctx.Value(context.UserIDContextKey).(string)
	items, err := s.storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	if len(items) == 0 {
		return nil, status.Error(codes.NotFound, "Resourse is empty")
	}

	response := &pb.GetURLByUserResponse{}
	for _, v := range items {
		responseElement := &pb.URLByUserResponseElement{
			ShortURL:    urlformat.FormatURL(string(s.config.BaseURL), v.Short),
			OriginalURL: v.Full,
		}
		response.Urls = append(response.Urls, responseElement)
	}

	return response, nil
}

// Обработка /short.ShortService/DeleteBatchURL
// DeleteBatchURL пакетно удаляет сокращенные ссылки
func (s *ShortGRPCService) DeleteBatchURL(ctx _context.Context, in *pb.DeleteBatchURLRequest) (*emptypb.Empty, error) {
	if _, isDeleter := s.storage.(storage.StorageDeleter); !isDeleter {
		return nil, status.Error(codes.Internal, "Batch delete is not supported for used storage type")
	}

	// пачки нет
	if len(in.Urls) == 0 {
		return nil, status.Error(codes.NotFound, "Resourse is empty")
	}
	s.storage.(storage.StorageDeleter).DeleteBatch(ctx, ctx.Value(context.UserIDContextKey).(string), []string(in.Urls))

	return nil, nil
}

// Обработка /short.ShortService/GetStats
// GetStats стастистика сокращенных URL и пользователей
func (s *ShortGRPCService) GetStats(ctx _context.Context, in *emptypb.Empty) (*pb.GetStatsResponse, error) {
	result, err := s.storage.GetStats(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.GetStatsResponse{
		Urls:  int32(result.URLs),
		Users: int32(result.Users),
	}, nil
}
