package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type Storage struct {
	engine EngineLayer
	logger *slog.Logger
}

var errInvalidLogger error = errors.New("invalid logger")

func NewStorage(engine EngineLayer, logger *slog.Logger) (*Storage, error) {
	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Storage{
		engine: engine,
		logger: logger,
	}, nil
}

type EngineLayer interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Del(ctx context.Context, key string) error
}

func (s Storage) Get(ctx context.Context, key string) (string, error) {
	logAttrs := []any{
		slog.String("component", "storage"),
		slog.String("method", "Get"),
		slog.String("key", key),
	}

	res, err := s.engine.Get(ctx, key)
	if err != nil {
		wErr := fmt.Errorf("get from engine: %w", err)
		s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return "", wErr
	}

	return res, nil
}

func (s Storage) Set(ctx context.Context, key, value string) error {
	logAttrs := []any{
		slog.String("component", "storage"),
		slog.String("method", "Set"),
		slog.String("key", key),
		slog.String("value", value),
	}

	err := s.engine.Set(ctx, key, value)
	if err != nil {
		wErr := fmt.Errorf("set to engine: %w", err)
		s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return wErr
	}

	return nil
}

func (s Storage) Del(ctx context.Context, key string) error {
	logAttrs := []any{
		slog.String("component", "storage"),
		slog.String("method", "Del"),
		slog.String("key", key),
	}

	err := s.engine.Del(ctx, key)
	if err != nil {
		wErr := fmt.Errorf("delete from engine: %w", err)
		s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return wErr
	}

	return nil
}
