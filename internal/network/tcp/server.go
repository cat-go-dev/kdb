package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
)

type Server struct {
	opts   *ServerOpts
	logger *slog.Logger
}

type ServerOpts struct {
	Port           uint
	ChannelBufSize int
}

const (
	defaultPort    = 8000
	defaultBufSize = 100
)

func NewServer(logger *slog.Logger, opts *ServerOpts) (*Server, error) {
	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Server{
		opts:   opts,
		logger: logger,
	}, nil
}

func (s Server) Run(ctx context.Context) error {
	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "Run"),
	}

	listener, err := net.Listen("tcp", s.getAddress())
	if err != nil {
		wErr := fmt.Errorf("trying to run tcp server: %w", err)
		s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return wErr
	}

	defer func() {
		err = listener.Close()
		if err != nil {
			s.logger.ErrorContext(ctx, fmt.Errorf("trying to close listener: %w", err).Error(), logAttrs...)
		}
	}()

	s.logger.InfoContext(ctx, "server is running", logAttrs...)

	for {
		select {
		case conn := <-s.getConnections(ctx, listener):
			go func() {
				defer func() {
					if r := recover(); r != nil {
						s.logger.ErrorContext(ctx, "caught panic: %s", r)
					}
				}()

				err = s.handleConnection(ctx, conn)
				if err != nil {
					s.logger.ErrorContext(ctx, fmt.Errorf("trying to handle connection: %w", err).Error(), logAttrs...)
				}
			}()
		case <-ctx.Done():
			s.logger.WarnContext(ctx, "server stopped by canceled context", logAttrs...)
			return errCanceledContext
		}
	}
}

func (s Server) getConnections(ctx context.Context, listener net.Listener) <-chan net.Conn {
	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "acceptConnection"),
	}

	connCh := make(chan net.Conn, s.getBufSize())

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.logger.ErrorContext(ctx, fmt.Errorf("tryint to accept connection: %w", err).Error(), logAttrs...)
				continue
			}

			connCh <- conn
		}
	}()

	return connCh
}

func (s Server) handleConnection(ctx context.Context, conn net.Conn) error {
	defer conn.Close()

	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "handleConnection"),
	}

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			wErr := fmt.Errorf("trying to read conn string: %w", err)
			s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
			return wErr
		}

		message = strings.TrimSpace(message)

		// todo: call database
		response := fmt.Sprintf("Got message: %v\n", message)
		_, err = conn.Write([]byte(response))
		if err != nil {
			wErr := fmt.Errorf("trying to response: %w", err)
			s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
			return wErr
		}
	}
}

func (s Server) getAddress() string {
	port := defaultPort
	if s.opts != nil && s.opts.Port != 0 {
		port = int(s.opts.Port)
	}

	return fmt.Sprintf(":%d", port)
}

func (s Server) getBufSize() int {
	size := defaultBufSize
	if s.opts != nil && s.opts.ChannelBufSize != 0 {
		size = int(s.opts.Port)
	}

	return size
}
