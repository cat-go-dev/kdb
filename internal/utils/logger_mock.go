package utils

import (
	"bytes"
	"log/slog"
)

func NewMockedLogger() *slog.Logger {
	buf := new(bytes.Buffer)
	return slog.New(slog.NewTextHandler(buf, nil))
}
