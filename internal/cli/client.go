package cli

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"kdb/internal/ports"
)

type Client struct {
	db     ports.Database
	logger *slog.Logger
}

func NewClient(db ports.Database, logger *slog.Logger) (*Client, error) {
	if db == nil {
		return nil, errInvalidDatabase
	}

	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Client{
		db:     db,
		logger: logger,
	}, nil
}

const commandPrefix = "[kdb] > "

func (c Client) Run(ctx context.Context) error {
	logAttrs := []any{
		slog.String("component", "cli"),
		slog.String("method", "Run"),
	}

	reader := bufio.NewReader(os.Stdin)
	commandCh := make(chan string)

	for {
		fmt.Print(commandPrefix)

		go func() {
			command, err := reader.ReadString('\n')
			if err != nil {
				wErr := fmt.Errorf("read command: %w", err)
				c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
				fmt.Printf("%ssomething went wrong \r\n", commandPrefix)
				return
			}

			preparedCommand := strings.ReplaceAll(command, "\r\n", "")

			commandCh <- preparedCommand
		}()

		select {
		case command := <-commandCh:
			fmt.Printf("%s%s \r\n", commandPrefix, c.executeCommand(ctx, command))
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

	r, err := c.db.Execute(ctx, command)
	if err != nil {
		wErr := fmt.Errorf("db executing: %w", err)
		c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return err.Error()
	}

	return r.Msg
}
