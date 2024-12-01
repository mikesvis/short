// Пакет перехватичков
package interceptors

import (
	_context "context"
	_errors "errors"
	"slices"
	"time"

	_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Регистрация по ключу jwt.AuthorizationMDKeyName. В результате успешной регистрации будет создан ключ и прописан ID пользователя в контекст.
func SignInInterceptor(ctx _context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	endPoints := []string{"/short.ShortService/SaveURL", "/short.ShortService/ShortenBatchURL"}
	if !slices.Contains(endPoints, info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	// Metadata есть
	if ok {
		tokenSlice, ok := md[jwt.AuthorizationMDKeyName]
		if ok && len(tokenSlice) > 0 {
			tokenString := tokenSlice[0]
			var userID string
			userID, err := jwt.GetUserIDFromTokenString(tokenString)

			// все ОК, пишем в контекст userID
			if err == nil {
				ctx := setUserIDToContext(ctx, userID)
				return handler(ctx, req)
			}

			// если пустой userID: Unauthenticated
			if _errors.Is(err, errors.ErrEmptyUserID) {
				return nil, status.Error(codes.Unauthenticated, "Authorization user id is empty")
			}

			// с токеном проблема, но не проблема подписи StatusInternalServerError
			if !_errors.Is(err, _jwt.ErrSignatureInvalid) {
				return nil, status.Error(codes.Internal, "Internal server error while processing auth token")
			}
		}
	}

	// Metadata нет или нет ключа или проблема подписи
	userID := uuid.NewString()
	expirationTime := time.Now().Add(jwt.TokenDuration)
	tokenString, err := jwt.CreateTokenString(userID, expirationTime)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error while creating auth token")
	}

	// записываем новый токен в metadata
	header := metadata.Pairs(jwt.AuthorizationMDKeyName, tokenString)
	grpc.SetHeader(ctx, header)

	// записываем UserID в контекст
	ctx = setUserIDToContext(ctx, userID)
	return handler(ctx, req)
}

// Авторизация по ключу jwt.AuthorizationMDKeyName
func AuthInterceptor(ctx _context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	endPoints := []string{"/short.ShortService/GetURLByUser", "/short.ShortService/DeleteBatchURL"}
	if !slices.Contains(endPoints, info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadata is empty")
	}

	tokenSlice, ok := md[jwt.AuthorizationMDKeyName]
	// Ключ не найден в Metadata
	if !ok || len(tokenSlice) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata token %s is not found", jwt.AuthorizationMDKeyName)
	}

	tokenString := tokenSlice[0]
	userID, err := jwt.GetUserIDFromTokenString(tokenString)

	// проблема с расшифровкой или валидностью JWT
	if err != nil && _errors.Is(err, errors.ErrInvalidToken) {
		return nil, status.Error(codes.Unauthenticated, "Failed to decode or token is invalid")
	}

	// UserID есть но он пустой
	if err != nil && _errors.Is(err, errors.ErrEmptyUserID) {
		return nil, status.Error(codes.Unauthenticated, "Empty User ID in token")
	}

	// какая-то другая проблема с токеном
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	// записываем UserID в контекст
	ctx = setUserIDToContext(ctx, userID)
	return handler(ctx, req)
}

func setUserIDToContext(ctx _context.Context, userID string) _context.Context {
	return _context.WithValue(ctx, context.UserIDContextKey, userID)
}
