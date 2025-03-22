package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"kdb/internal/ports"
)

type Server struct {
	opts     *ServerOpts
	logger   *slog.Logger
	executor Executor
	conns    map[string]struct{}
}

type ServerOpts struct {
	Host           string
	Port           uint
	MaxConnections uint
}

type Executor interface {
	Execute(ctx context.Context, commandStr string) (*ports.Result, error)
}

const (
	defaultHost           = "localhost"
	defaultPort           = 8000
	defaultBufSize        = 100
	defaultMaxConnections = 100
)

func NewServer(executor Executor, logger *slog.Logger, opts *ServerOpts) (*Server, error) {
	if executor == nil {
		return nil, errInvalidExecutor
	}

	if logger == nil {
		return nil, errInvalidLogger
	}

	if opts.MaxConnections == 0 {
		opts.MaxConnections = defaultMaxConnections
	}

	return &Server{
		executor: executor,
		opts:     opts,
		logger:   logger,
		conns: make(map[string]struct{}, opts.MaxConnections),
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "Run"),
		slog.String("address", s.getAddress()),
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

	conns := s.getConnections(ctx, listener)

	for {
		select {
		case conn := <-conns:
			go func() {
				defer func() {
					if r := recover(); r != nil {
						s.logger.ErrorContext(ctx, "caught panic: %s", r)
					}
				}()
				defer func() {
					s.logger.InfoContext(ctx, fmt.Sprintf("conn %v is closed", conn.RemoteAddr().String()), logAttrs...)
					delete(s.conns, conn.RemoteAddr().String())
					conn.Close()
				}()

				if len(s.conns) == int(s.opts.MaxConnections) {
					logAttrs = append(logAttrs, slog.String("remote_addr", conn.RemoteAddr().String()))
					s.logger.WarnContext(ctx, "connection limit reached", logAttrs...)

					err := s.rejectMaxConnCount(ctx, conn)
					if err != nil {
						s.logger.ErrorContext(ctx, fmt.Errorf("trying to reject connection: %w", err).Error(), logAttrs...)
					}

					return
				}

				s.conns[conn.RemoteAddr().String()] = struct{}{}

				s.logger.InfoContext(ctx, fmt.Sprintf("new conn %v", conn.RemoteAddr().String()), logAttrs...)

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

func (s *Server) getConnections(ctx context.Context, listener net.Listener) <-chan net.Conn {
	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "acceptConnection"),
	}

	connCh := make(chan net.Conn, defaultBufSize)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				wErr := fmt.Errorf("%s: %w", errTryingToAcceptConnection, err)
				s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
				continue
			}

			connCh <- conn
		}
	}()

	return connCh
}

func (s Server) handleConnection(ctx context.Context, conn net.Conn) error {
	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "handleConnection"),
	}

	reader := bufio.NewReader(conn)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			wErr := fmt.Errorf("trying to read conn string: %w", err)
			s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
			return wErr
		}

		command = strings.TrimSpace(command)
		s.logger.InfoContext(ctx, fmt.Sprintf("Got message: %v", command), logAttrs...)

		var response string
		result, err := s.executor.Execute(ctx, command)
		if err != nil {
			wErr := fmt.Errorf("execute error: %w", err)
			s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
			response = fmt.Sprintf("An error while executing command: %s", command)
		} else {
			response = result.Msg
		}

		response = fmt.Sprintf("%s\n", response)

		_, err = conn.Write([]byte(response))
		if err != nil {
			wErr := fmt.Errorf("trying to response: %w", err)
			s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
			return wErr
		}
	}
}

func (s Server) rejectMaxConnCount(ctx context.Context, conn net.Conn) error {
	logAttrs := []any{
		slog.String("component", "tcp_server"),
		slog.String("method", "rejectMaxConnCount"),
	}

	response := fmt.Sprintf("max connection limit reached\n")

	_, err := conn.Write([]byte(response))
	if err != nil {
		wErr := fmt.Errorf("trying to response: %w", err)
		s.logger.ErrorContext(ctx, wErr.Error(), logAttrs...)
		return wErr
	}

	return nil
}

func (s Server) getAddress() string {
	host := defaultHost
	if s.opts != nil && s.opts.Host != "" {
		host = s.opts.Host
	}

	port := defaultPort
	if s.opts != nil && s.opts.Port != 0 {
		port = int(s.opts.Port)
	}

	return fmt.Sprintf("%s:%d", host, port)
}
