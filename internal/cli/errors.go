package cli

import "errors"

var (
	errInvalidDatabase = errors.New("invalid database")
	errInvalidLogger   = errors.New("invalid logger")
	errCanceledContext = errors.New("canceled context")
)
