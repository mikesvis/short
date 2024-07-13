package errors

import ge "errors"

// Чуть не сошел сума с циклической зависимостью :(
var ErrConflict = ge.New("conflict error")
