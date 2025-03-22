package database

import "errors"

var (
	errInvalidLogger  = errors.New("invalid logger")
	errInvalidCompute = errors.New("invalid compute")
	errInvalidStorage = errors.New("invalid storage")
	errUnknownCommand = errors.New("unknown command")
	errComputeParse   = errors.New("compute parse")
)
