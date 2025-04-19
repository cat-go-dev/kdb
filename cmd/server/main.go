package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"kdb/internal/config"
	"kdb/internal/database"
	"kdb/internal/database/compute"
	"kdb/internal/database/storage"
	"kdb/internal/database/storage/engine"
	logger "kdb/internal/logs"
	"kdb/internal/network/tcp"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

const defaultEnv string = "local"

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	env, ok := os.LookupEnv("ENV")
	if !ok {
		fmt.Println(ctx, "empty env enviroment var")
		env = defaultEnv
	}

	cfg, err := config.NewConfig()
	if err != nil {
		wErr := fmt.Errorf("creating config: %w", err)
		fmt.Printf("new config error: %s\n", wErr.Error())
		return
	}
	err = cfg.Init(ctx, env)
	if err != nil {
		wErr := fmt.Errorf("config init: %w", err)
		fmt.Printf("config init error: %s\n", wErr.Error())
		return
	}

	logger := logger.NewLogger(ctx, *cfg.Data)

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

	tcpServer, err := tcp.NewServer(database, logger, &tcp.ServerOpts{
		Host:           cfg.Data.Network.Host,
		Port:           uint(cfg.Data.Network.Port),
		MaxConnections: uint(cfg.Data.Network.MaxConnections),
	})
	if err != nil {
		wErr := fmt.Errorf("creating tcp server: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	err = tcpServer.Run(ctx)
	if err != nil {
		wErr := fmt.Errorf("running tcp server: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
	}

	cancel()
}
