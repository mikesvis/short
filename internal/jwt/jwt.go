// Пакет обработки jwt.
package jwt

import (
	"time"

	_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mikesvis/short/internal/errors"
)

// SecretPass используется для шифровки/дешифровки.
const SecretPass = "mySecretPass"

// Имя куки авторизации.
const AuthorizationCookieName = "Authorization-JWT"

// Ключ авторизации в Metadata
const AuthorizationMDKeyName = "authorization-jwt"

// Время жизни куки авторизации.
const TokenDuration = time.Hour * 24 * 30

// Claims в JWT.
type Claims struct {
	// ID пользователя.
	UserID string `json:"userId"`
	_jwt.RegisteredClaims
}

// Получение ID пользователя из токена.
func GetUserIDFromTokenString(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := _jwt.ParseWithClaims(tokenString, claims, func(token *_jwt.Token) (any, error) {
		return []byte(SecretPass), nil
	})

	// все хорошо в куке, не трогаем
	if err == nil && token.Valid {
		// пустой UserID в токене (по заданию)
		if len(claims.UserID) == 0 {
			return "", errors.ErrEmptyUserID
		}

		return claims.UserID, nil
	}

	if err == nil && !token.Valid {
		return "", errors.ErrInvalidToken
	}

	return "", err
}

// Создание токена авторизации.
func CreateTokenString(userID string, exp time.Time) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: _jwt.RegisteredClaims{
			ExpiresAt: _jwt.NewNumericDate(exp),
		},
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretPass))
	if err != nil {
		return "", err
	}

	return tokenString, err
}
