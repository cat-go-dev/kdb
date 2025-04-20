package wal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"kdb/internal/config"
)

// todo: change to 100
const (
	defaultInChanSize = 2
	defaultChunkSize  = 3
)

type WAL struct {
	in       chan string
	chunk    []string
	fWritter fWritter
	logger   *slog.Logger
}

func NewWAL(config *config.AppConfig, logger *slog.Logger) (*WAL, error) {
	if config == nil {
		return nil, errInvalidConfig
	}

	if logger == nil {
		return nil, errInvalidLogger
	}

	writter, err := newWritter(config.Data.WAL.DataDitrectory, logger)
	if err != nil {
		return nil, fmt.Errorf("creating writter error: %w", err)
	}

	return &WAL{
		in:       make(chan string, defaultInChanSize),
		chunk:    make([]string, 0, defaultChunkSize),
		fWritter: *writter,
		logger:   logger,
	}, nil
}

func (w *WAL) Run(ctx context.Context) {
	go w.listenAndWrite(ctx)
}

func (w *WAL) Write(ctx context.Context, log string) {
	// maybe need block
	w.in <- log
}

func (w *WAL) listenAndWrite(ctx context.Context) {
	logAttrs := []any{
		slog.String("component", "WAL"),
		slog.String("method", "listenAndWrite"),
	}

	for {
		select {
		case log := <-w.in:
			w.handle(ctx, log, false)
		case <-ctx.Done():
			w.logger.WarnContext(ctx, "wal stopped by canceled context", logAttrs...)
			return
		}
	}
}

func (w *WAL) handle(ctx context.Context, log string, force bool) {
	logAttrs := []any{
		slog.String("component", "WAL"),
		slog.String("method", "handle"),
		slog.Bool("force", force),
	}

	w.chunk = append(w.chunk, w.prepareLog(log))
	if !force && len(w.chunk) < defaultChunkSize {
		return
	}

	err := w.fWritter.write(ctx, w.chunk)
	if err != nil {
		w.logger.ErrorContext(ctx, fmt.Errorf("writter: %w", err).Error(), logAttrs...)
	}

	w.chunk = w.chunk[:0]
	go w.logger.InfoContext(ctx, "logs written", logAttrs...)
}

func (w *WAL) prepareLog(log string) string {
	return fmt.Sprintf("%d %s", time.Now().Unix(), log)
}
