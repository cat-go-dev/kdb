package compute

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommandType(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected CommandType
	}{
		{name: "should be get command type", command: "GET", expected: Get},
		{name: "should be set command type", command: "SET", expected: Set},
		{name: "should be del command type", command: "DEL", expected: Del},
	}

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandType, _ := compute.getCommandType(tt.command)
			assert.Equal(t, tt.expected, commandType)
		})
	}
}

func TestGetArgumentsSuccess(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected Arguments
	}{
		{name: "should be only key argument", tokens: []string{"GET", "test"}, expected: Arguments{Key: "test"}},
		{name: "should be key and value argument", tokens: []string{"GET", "test", "true"}, expected: Arguments{Key: "test", Value: "true"}},
	}

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := compute.getArguments(tt.tokens)

			assert.Nil(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetArgumentsNotEnoghtArguments(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected error
	}{
		{name: "1 argument error", tokens: []string{"GET"}, expected: errNotEnoughArguments},
		{name: "0 argument error", tokens: []string{}, expected: errNotEnoughArguments},
		{name: "2 argument without error", tokens: []string{"GET", "test"}, expected: nil},
	}

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := compute.getArguments(tt.tokens)
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestParseGetCommand(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	commandType, _ := compute.getCommandType("GET")
	expected := &Command{
		Type: commandType,
		Arguments: Arguments{
			Key: Argument("test"),
		},
	}

	commandStr := "GET test"

	actual, err := compute.Parse(ctx, commandStr)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseSetCommand(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	commandType, _ := compute.getCommandType("SET")
	expected := &Command{
		Type: commandType,
		Arguments: Arguments{
			Key:   Argument("test"),
			Value: Argument("true"),
		},
	}

	commandStr := "SET test true"

	actual, err := compute.Parse(ctx, commandStr)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseDelCommand(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	commandType, _ := compute.getCommandType("DEL")
	expected := &Command{
		Type: commandType,
		Arguments: Arguments{
			Key: Argument("test"),
		},
	}

	commandStr := "DEL test"

	actual, err := compute.Parse(ctx, commandStr)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseNotEnoughArguments(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	commandStr := "GET"

	expected := errNotEnoughArguments

	actual, err := compute.Parse(ctx, commandStr)

	assert.Nil(t, actual)
	assert.NotNil(t, err)
	assert.Equal(t, expected, err)
}

func TestParseUnknownCommandType(t *testing.T) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	compute, _ := NewCompute(slog.New(slog.NewTextHandler(buf, nil)))

	commandStr := "TEST test"

	expected := errUnknownCommandType

	actual, err := compute.Parse(ctx, commandStr)

	assert.Nil(t, actual)
	assert.NotNil(t, err)
	assert.Equal(t, expected, err)
}
