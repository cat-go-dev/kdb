package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"kdb/internal/config"
)

const (
	levelDebug   = "debug"
	levelInfo    = "info"
	levelWarning = "warning"
	levelError   = "error"
)

func NewLogger(ctx context.Context, cfg config.Config) *slog.Logger {
	logFile, err := os.OpenFile(getFileName(cfg.Logging.OutputDir), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Can't open log file: %v", err))
	}

	go func() {
		<-ctx.Done()
		logFile.Close()
	}()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: getLevel(cfg.Logging.Level),
	}))
}

func getLevel(level string) slog.Level {
	switch level {
	case levelDebug:
		return slog.LevelDebug
	case levelInfo:
		return slog.LevelInfo
	case levelWarning:
		return slog.LevelWarn
	case levelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}


func getFileName(dir string) string {
	t := time.Now()
	return fmt.Sprintf("%s/log_%s.log", dir, t.Format(time.DateOnly))
}
