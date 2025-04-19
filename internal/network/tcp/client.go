package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

type Client struct {
	in     chan string
	out    chan string
	logger *slog.Logger
	opts   *ClientOpts
}

type ClientOpts struct {
	Server  string
	Port    int
}

const (
	defaultServer     = "localhost"
	defaultServerPort = "8000"
)

func NewClient(logger *slog.Logger, opts *ClientOpts) (*Client, error) {
	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Client{
		in:     make(chan string),
		out:    make(chan string),
		logger: logger,
		opts:   opts,
	}, nil
}

func (c *Client) Run(ctx context.Context) error {
	logAttrs := []any{
		slog.String("component", "tcp_client"),
		slog.String("method", "Run"),
	}

	conn, err := net.Dial("tcp", c.getAddressToConnect())
	if err != nil {
		wErr := fmt.Errorf("trying create tcp client connection: %w", err)
		c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return wErr
	}

	c.logger.InfoContext(ctx, "tcp client is started", logAttrs...)

	go func() {
		defer conn.Close()
		for {
			select {
			case message := <-c.in:
				message = strings.TrimSpace(message)

				_, err := conn.Write([]byte(message + "\n"))
				if err != nil {
					wErr := fmt.Errorf("trying to send message to server: %w", err)
					c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
					c.out <- "internal error"
					continue
				}

				response, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					wErr := fmt.Errorf("trying to read response from server: %w", err)
					c.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
					c.out <- "internal error"
					continue
				}

				c.out <- response
			case <-ctx.Done():
				c.logger.WarnContext(ctx, "client stopped by canceled context", logAttrs...)
				return
			}

		}
	}()

	return nil
}

func (c *Client) Call(ctx context.Context, message string) (string, error) {
	logAttrs := []any{
		slog.String("component", "tcp_client"),
		slog.String("method", "Call"),
	}

	select {
	case c.in <- message:
	case <-ctx.Done():
		c.logger.WarnContext(ctx, "client stopped by canceled context", logAttrs...)
		return "", errCanceledContext
	}

	select {
	case message := <-c.out:
		return message, nil
	case <-ctx.Done():
		c.logger.WarnContext(ctx, "client stopped by canceled context", logAttrs...)
		return "", errCanceledContext
	}
}

func (c *Client) getAddressToConnect() string {
	server := defaultServer
	if c.opts.Server != "" {
		server = c.opts.Server
	}

	port := defaultPort
	if c.opts.Port != 0 {
		port = c.opts.Port
	}

	return fmt.Sprintf("%s:%d", server, port)
}
