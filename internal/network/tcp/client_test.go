package tcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"kdb/internal/utils"
)

func TestNewClientEmptyLogger(t *testing.T) {
	_, err := NewClient(nil, nil)
	assert.ErrorIs(t, err, errInvalidLogger)
}

func TestNewClientSuccess(t *testing.T) {
	logger := utils.NewMockedLogger()

	_, err := NewClient(logger, nil)
	assert.NoError(t, err)
}

func TestRunWithInvalidAddress(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logger := utils.NewMockedLogger()

	client, err := NewClient(logger, &ClientOpts{
		Server: "123.123.123.123",
		Port:   21,
	})
	assert.Nil(t, err)

	err = client.Run(ctx)
	assert.Error(t, err)

	cancel()
}
