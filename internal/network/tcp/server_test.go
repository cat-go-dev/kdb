package tcp

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"kdb/internal/network/tcp/mocks"
)

func TestNewServerEmptyExecutor(t *testing.T) {
	_, err := NewServer(nil, nil, &ServerOpts{})
	assert.ErrorIs(t, err, errInvalidExecutor)
}

func TestNewServerEmptyLogger(t *testing.T) {
	executor := mocks.NewExecutor(t)
	_, err := NewServer(executor, nil, &ServerOpts{})
	assert.ErrorIs(t, err, errInvalidLogger)
}

func TestNewServerSuccess(t *testing.T) {
	executor := mocks.NewExecutor(t)
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(buf, nil))

	_, err := NewServer(executor, logger, &ServerOpts{})
	assert.NoError(t, err)
}

func TestRunServerInvalidAddress(t *testing.T) {
	ctx := context.Background()

	executor := mocks.NewExecutor(t)
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(buf, nil))

	expectedErr := errTryingToRunServer.Error()

	server, err := NewServer(executor, logger, &ServerOpts{
		Host: "test",
		Port: 123,
	})
	assert.NoError(t, err)

	err = server.Run(ctx)
	assert.ErrorContains(t, err, expectedErr)
}

func TestExitWithContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)

	executor := mocks.NewExecutor(t)
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(buf, nil))

	expectedErr := errTryingToRunServer.Error()

	server, err := NewServer(executor, logger, &ServerOpts{
		Host: "localhost",
		Port: 8000,
	})
	assert.NoError(t, err)

	err = server.Run(ctx)
	assert.ErrorContains(t, err, expectedErr)
	cancel()
}
