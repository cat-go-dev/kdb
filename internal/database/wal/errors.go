package wal

import "errors"

var (
	errInvalidConfig     = errors.New("empty config")
	errInvalidLogger     = errors.New("empty logger")
	errEmptyWALDirectory = errors.New("empty wal directory")
)
