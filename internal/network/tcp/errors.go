package tcp

import "errors"

var (
	errInvalidLogger            = errors.New("invalid logger")
	errInvalidExecutor          = errors.New("invalid executor")
	errCanceledContext          = errors.New("canceled context")
	errTryingToRunServer        = errors.New("trying to run tcp server")
	errTryingToAcceptConnection = errors.New("tryint to accept connection")
)
