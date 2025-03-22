package main

import (
	"context"
	"fmt"
	"kdb/internal/cli"
	"kdb/internal/database"
	"kdb/internal/database/compute"
	"kdb/internal/database/storage"
	"kdb/internal/database/storage/engine"
	"kdb/internal/network/tcp"
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

	tcpServer, err := tcp.NewServer(logger, &tcp.ServerOpts{
		Port:           8000, // todo: get from config
		ChannelBufSize: 1000, // todo: get from config
	})
	if err != nil {
		wErr := fmt.Errorf("creating tcp server: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	go func() {
		err = tcpServer.Run(ctx)
		if err != nil {
			wErr := fmt.Errorf("running tcp server: %w", err)
			logger.ErrorContext(ctx, wErr.Error())
			return
		}
	}()

	err = client.Run(ctx)
	if err != nil {
		wErr := fmt.Errorf("running cli client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}
}
