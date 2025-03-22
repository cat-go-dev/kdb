package tcp

import "errors"

var (
	errInvalidLogger   = errors.New("invalid logger")
	errCanceledContext = errors.New("canceled context")
)
