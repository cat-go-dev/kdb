package cli

import "errors"

var (
	errInvalidTcpClient = errors.New("invalid tcp client")
	errInvalidLogger    = errors.New("invalid logger")
	errCanceledContext  = errors.New("canceled context")
)
