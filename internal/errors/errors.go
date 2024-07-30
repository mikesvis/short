package errors

import _goerrors "errors"

// Чуть не сошел сума с циклической зависимостью :(
var ErrConflict = _goerrors.New("conflict error")
var ErrEmptyUserID = _goerrors.New("empty user id")
var ErrInvalidToken = _goerrors.New("invalid token")
