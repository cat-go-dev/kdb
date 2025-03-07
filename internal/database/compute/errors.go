package compute

import "errors"

var (
	errInvalidLogger      = errors.New("invalid logger")
	errNotEnoughArguments = errors.New("not enought arguments")
	errUnknownCommandType = errors.New("unknown command type")
)
