package cli

import (
	"kdb/internal/network/tcp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientWithEmptyTcpClient(t *testing.T) {
	client, err := NewClient(nil, nil)

	assert.Nil(t, client)
	assert.ErrorIs(t, err, errInvalidTcpClient)
}

func TestNewClientWithEmptyLogger(t *testing.T) {
	tcpClient := tcp.Client{}
	client, err := NewClient(&tcpClient, nil)

	assert.Nil(t, client)
	assert.ErrorIs(t, err, errInvalidLogger)
}
