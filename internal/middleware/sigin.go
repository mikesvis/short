package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mikesvis/short/internal/domain"
)

const secretPass = "mySecretPass"
const AuthorizationCookieName = "Authorization-JWT"
const tokenDuration = time.Hour * 24 * 30

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

func SignIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(AuthorizationCookieName)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// кука есть
		if err == nil {
			tokenString := authCookie.Value
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				return []byte(secretPass), nil
			})

			// все хорошо в куке, не трогаем
			if err == nil && token.Valid {
				// пустой UserID в токене (по заданию)
				if len(claims.UserID) == 0 {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				ctx := setUserIDToContext(r, claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// с кукой что-то не так, но это не проблема подписи
			if err != nil && !errors.Is(err, jwt.ErrSignatureInvalid) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// куки нет или проблема подписи - создаем новую
		userID := uuid.NewString()
		expirationTime := time.Now().Add(tokenDuration)
		tokenString, err := createTokenString(userID, expirationTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    AuthorizationCookieName,
			Value:   tokenString,
			Expires: expirationTime,
		})

		ctx := setUserIDToContext(r, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func createTokenString(userID string, exp time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretPass))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func setUserIDToContext(r *http.Request, userID string) context.Context {
	return context.WithValue(r.Context(), domain.ContextUserKey, userID)
}
