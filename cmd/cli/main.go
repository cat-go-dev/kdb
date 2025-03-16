package main

import (
	"context"
	"fmt"
	"kdb/internal/cli"
	"kdb/internal/database"
	"kdb/internal/database/compute"
	"kdb/internal/database/storage"
	"kdb/internal/database/storage/engine"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	// todo: maybe simple logs (without JSON) for localhost
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	compute, err := compute.NewCompute(logger)
	if err != nil {
		wErr := fmt.Errorf("creating compute: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	storage, err := storage.NewStorage(engine.NewEngine(), logger)
	if err != nil {
		wErr := fmt.Errorf("creating storage: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	database, err := database.NewDatabase(compute, storage, logger)
	if err != nil {
		wErr := fmt.Errorf("creating database: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	client, err := cli.NewClient(database, logger)
	if err != nil {
		wErr := fmt.Errorf("creating cli client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	err = client.Run(ctx)
	if err != nil {
		wErr := fmt.Errorf("running cli client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}
}
