package database

import "errors"

var (
	errInvalidLogger  = errors.New("invalid logger")
	errInvalidCompute = errors.New("invalid compute")
)
