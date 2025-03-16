package database

import (
	"context"
	"fmt"
	"log/slog"

	"kdb/internal/database/compute"
	"kdb/internal/ports"
)

type Database struct {
	compute *compute.Compute
	storage StorageLayer
	logger  *slog.Logger
}

type StorageLayer interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Del(ctx context.Context, key string) error
}

func NewDatabase(compute *compute.Compute, storage StorageLayer, logger *slog.Logger) (*Database, error) {
	if compute == nil {
		return nil, errInvalidCompute
	}

	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Database{
		compute: compute,
		storage: storage,
		logger:  logger,
	}, nil
}

func (d Database) Execute(ctx context.Context, commandStr string) (*ports.Result, error) {
	logAttrs := []any{
		slog.String("component", "database"),
		slog.String("method", "Execute"),
	}

	command, err := d.compute.Parse(ctx, commandStr)
	if err != nil {
		wErr := fmt.Errorf("compute parse: %w", err)
		d.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return nil, err
	}

	return d.executeCommand(ctx, command)
}

// todo: make test
func (d Database) executeCommand(ctx context.Context, command *compute.Command) (*ports.Result, error) {
	logAttrs := []any{
		slog.String("component", "database"),
		slog.String("method", "executeCommand"),
		slog.Any("command", command),
	}

	var res string
	var err error

	switch {
	case command.Type.IsGet():
		res, err = d.storage.Get(ctx, string(command.Arguments.Key))
		logAttrs = append(logAttrs, slog.String("storage method", "get"))
	case command.Type.IsSet():
		err = d.storage.Set(ctx, string(command.Arguments.Key), string(command.Arguments.Value))
		logAttrs = append(logAttrs, slog.String("storage method", "set"))
	case command.Type.IsDel():
		err = d.storage.Del(ctx, string(command.Arguments.Key))
		logAttrs = append(logAttrs, slog.String("storage method", "del"))
	}

	if err != nil {
		wErr := fmt.Errorf("storage call: %w", err)
		d.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)

		return nil, wErr
	}

	return &ports.Result{
		Msg: res,
	}, nil
}
