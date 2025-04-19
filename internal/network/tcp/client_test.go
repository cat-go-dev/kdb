package tcp

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientEmptyLogger(t *testing.T) {
	_, err := NewClient(nil, nil)
	assert.ErrorIs(t, err, errInvalidLogger)
}

func TestNewClientSuccess(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(buf, nil))

	_, err := NewClient(logger, nil)
	assert.NoError(t, err)
}

func TestRunWithInvalidAddress(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	buf := new(bytes.Buffer)
	logger := slog.New(slog.NewTextHandler(buf, nil))

	client, err := NewClient(logger, &ClientOpts{
		Server: "123.123.123.123",
		Port:   21,
	})
	assert.Nil(t, err)

	err = client.Run(ctx)
	assert.Error(t, err)

	cancel()
}
