// Модуль авторизации в приложении.
package middleware

import (
	_context "context"
	_errors "errors"
	"net/http"
	"time"

	_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/context"
	"github.com/mikesvis/short/internal/errors"
	"github.com/mikesvis/short/internal/jwt"
)

// Регистрация по куке jwt.AuthorizationCookieName. В результате успешной регистрации будет создана кука и прописан ID пользователя в контекст.
func SignIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(jwt.AuthorizationCookieName)
		if err != nil && !_errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// кука есть
		if err == nil {
			tokenString := authCookie.Value

			userID, err := jwt.GetUserIDFromTokenString(tokenString)

			// все ОК, пишем в контекст userID
			if err == nil {
				ctx := setUserIDToContext(r, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// если пустой userID: StatusUnauthorized
			if _errors.Is(err, errors.ErrEmptyUserID) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// с токеном проблем, но не проблема подписи StatusInternalServerError
			if !_errors.Is(err, _jwt.ErrSignatureInvalid) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// куки нет или проблема подписи - создаем новую
		userID := uuid.NewString()
		expirationTime := time.Now().Add(jwt.TokenDuration)
		tokenString, err := jwt.CreateTokenString(userID, expirationTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, CreateAuthCookie(tokenString, expirationTime))

		ctx := setUserIDToContext(r, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Построение авторизационной куки.
func CreateAuthCookie(tokenString string, exp time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    jwt.AuthorizationCookieName,
		Value:   tokenString,
		Expires: exp,
		Path:    "/",
	}
}

// Авторизация по куке jwt.AuthorizationCookieName
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(jwt.AuthorizationCookieName)
		// проблема с получением куки или ее нет
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString := authCookie.Value
		userID, err := jwt.GetUserIDFromTokenString(tokenString)

		// проблема с расшифровкой или валидностью JWT
		if err != nil && _errors.Is(err, errors.ErrInvalidToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// UserID есть но он пустой
		if err != nil && _errors.Is(err, errors.ErrEmptyUserID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// какая-то другая проблема с токеном
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := setUserIDToContext(r, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setUserIDToContext(r *http.Request, userID string) _context.Context {
	return _context.WithValue(r.Context(), context.UserIDContextKey, userID)
}
