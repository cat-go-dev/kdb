package cli

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	"kdb/internal/database"
	"kdb/internal/ports"
	"kdb/internal/ports/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewClientWithEmptyDatabase(t *testing.T) {
	client, err := NewClient(nil, nil)

	assert.Nil(t, client)
	assert.ErrorIs(t, err, errInvalidDatabase)
}

func TestNewClientWithEmptyLogger(t *testing.T) {
	db := &database.Database{}
	client, err := NewClient(db, nil)

	assert.Nil(t, client)
	assert.ErrorIs(t, err, errInvalidLogger)
}

func TestExecuteCommandSuccess(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDatabase(t)
	command := "test"
	expected := "success"

	db.EXPECT().Execute(ctx, command).Return(&ports.Result{
		Msg: expected,
	}, nil).Once()

	buf := new(bytes.Buffer)
	client, err := NewClient(db, slog.New(slog.NewTextHandler(buf, nil)))

	assert.Nil(t, err)

	actual := client.executeCommand(ctx, command)

	assert.Equal(t, expected, actual)
}

func TestExecuteCommandError(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDatabase(t)
	command := "test"
	expected := "test"

	db.EXPECT().Execute(ctx, command).Return(nil, fmt.Errorf("test")).Once()

	buf := new(bytes.Buffer)
	client, err := NewClient(db, slog.New(slog.NewTextHandler(buf, nil)))

	assert.Nil(t, err)

	actual := client.executeCommand(ctx, command)

	assert.Equal(t, expected, actual)
}
