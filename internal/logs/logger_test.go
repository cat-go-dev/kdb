package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLevel(t *testing.T) {
	tests := []struct{
		name string
		level string
		slogLevel slog.Level
	}{
		{name: "debug level", level: levelDebug, slogLevel: slog.LevelDebug},
		{name: "info level", level: levelInfo, slogLevel: slog.LevelInfo},
		{name: "warning level", level: levelWarning, slogLevel: slog.LevelWarn},
		{name: "error level", level: levelError, slogLevel: slog.LevelError},
		{name: "unknown level", level: "notice", slogLevel: slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			assert.Equal(t, tt.slogLevel, getLevel(tt.level))
		})
	}
}
