package errors

import _goerrors "errors"

// Чуть не сошел сума с циклической зависимостью :(
var ErrConflict = _goerrors.New("conflict error")
