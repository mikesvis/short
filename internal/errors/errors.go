// Модуль ошибок в приложении.
package errors

import _goerrors "errors"

// Конфликт при сохранении.
var ErrConflict = _goerrors.New("conflict error")

// Не указан ID пользователя.
var ErrEmptyUserID = _goerrors.New("empty user id")

// Неправильный токен
var ErrInvalidToken = _goerrors.New("invalid token")
