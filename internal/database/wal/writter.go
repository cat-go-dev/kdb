package wal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type fWritter struct {
	fName   string
	logger *slog.Logger
}

const permissions = 0644

func newWritter(fName string, logger *slog.Logger) (*fWritter, error) {
	if fName == "" {
		return nil, errEmptyWALDirectory
	}

	if logger == nil {
		return nil, errInvalidLogger
	}

	return &fWritter{
		fName:  fName,
		logger: logger,
	}, nil
}

func (f *fWritter) write(ctx context.Context, rows []string) error {
	logAttrs := []any{
		slog.String("component", "fWritter"),
		slog.String("method", "write"),
	}

	file, err := os.OpenFile(f.fName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, permissions)
	if err != nil {
		wErr := fmt.Errorf("open file: %w", err)
		f.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(f.buildMessage(rows)); err != nil {
		wErr := fmt.Errorf("write string to file: %w", err)
		f.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return err
    }

	return nil
}

func (f *fWritter) buildMessage(rows []string) string {
	var message string

	for _, row := range rows {
		message = fmt.Sprintf("%s%s\n", message, row)
	}

	return message
}
