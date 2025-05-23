package main

import (
	"context"
	"fmt"
	"os"

	"kdb/internal/cli"
	"kdb/internal/config"
	logger "kdb/internal/logs"
	"kdb/internal/network/tcp"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
        panic("No .env file found")
    }
}

const defaultEnv string = "local"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// todo: maybe simple logs (without JSON) for localhost

	env, ok := os.LookupEnv("ENV")
	if !ok {
		fmt.Println("empty env enviroment var")
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

	tcpClient, err := tcp.NewClient(logger, &tcp.ClientOpts{
		Server:  cfg.Data.Network.Host,
		Port:    cfg.Data.Network.Port,
	})
	if err != nil {
		wErr := fmt.Errorf("creating tcp client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	err = tcpClient.Run(ctx)
	if err != nil {
		wErr := fmt.Errorf("running tcp client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	cli, err := cli.NewClient(tcpClient, logger)
	if err != nil {
		wErr := fmt.Errorf("creating cli client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
		return
	}

	err = cli.Run(ctx)
	if err != nil {
		wErr := fmt.Errorf("running cli client: %w", err)
		logger.ErrorContext(ctx, wErr.Error())
	}

	cancel()
}
