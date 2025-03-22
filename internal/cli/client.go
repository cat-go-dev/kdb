package cli

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Client struct {
	executor Executor
	logger   *slog.Logger
}

type Executor interface {
	Call(ctx context.Context, command string) (string, error)
}

func NewClient(executor Executor, logger *slog.Logger) (*Client, error) {
	if executor == nil {
		return nil, errInvalidTcpClient
	}

	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Client{
		executor: executor,
		logger:   logger,
	}, nil
}

const (
	commandPrefix = "[kdb] > "
	commandExit   = "exit"
)

func (c Client) Run(ctx context.Context) error {
	logAttrs := []any{
		slog.String("component", "cli"),
		slog.String("method", "Run"),
	}

	ctx, cancel := context.WithCancel(ctx)

	reader := bufio.NewReader(os.Stdin)

	commandCh := make(chan string)
	exitCh := make(chan struct{})
	defer func() {
		close(commandCh)
		close(exitCh)
		cancel()
	}()

	for {
		fmt.Print(commandPrefix)

		go func() {
			command, err := reader.ReadString('\n')
			if err != nil {
				wErr := fmt.Errorf("read command: %w", err)
				logAttrs = append(logAttrs, slog.String("raw_command", command))
				c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
				fmt.Printf("%ssomething went wrong \r\n", commandPrefix)
				return
			}

			preparedCommand := strings.ReplaceAll(command, "\r", "")
			preparedCommand = strings.ReplaceAll(preparedCommand, "\n", "")

			if preparedCommand == commandExit {
				exitCh <- struct{}{}
				return
			}

			commandCh <- preparedCommand
		}()

		select {
		case command := <-commandCh:
			fmt.Printf("%s%s \r\n", commandPrefix, c.executeCommand(ctx, command))
		case <-exitCh:
			c.logger.InfoContext(ctx, "Exit command", logAttrs...)
			return nil
		case <-ctx.Done():
			c.logger.WarnContext(ctx, "canceled context", logAttrs...)
			return errCanceledContext
		}
	}
}

func (c Client) executeCommand(ctx context.Context, command string) string {
	logAttrs := []any{
		slog.String("component", "cli"),
		slog.String("method", "executeCimmand"),
	}

	r, err := c.executor.Call(ctx, command)
	if err != nil {
		wErr := fmt.Errorf("db executing: %w", err)
		c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return err.Error()
	}

	return r
}
