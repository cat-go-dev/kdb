package database

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"kdb/internal/database/compute"
	"kdb/internal/database/mocks"
)

func TestNewDatabaseEmptyCompute(t *testing.T) {
	expectedErr := errInvalidCompute

	_, err := NewDatabase(nil, nil, nil)

	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestNewDatabaseEmptyStorage(t *testing.T) {
	compute := getMockedCompute(t)

	expectedErr := errInvalidStorage

	_, err := NewDatabase(compute, nil, nil)

	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestNewDatabaseEmptyLogger(t *testing.T) {
	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)

	expectedErr := errInvalidLogger

	_, err := NewDatabase(compute, storage, nil)

	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestNotEnoughArguments(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := getMockedLogger()

	db, err := NewDatabase(compute, storage, logger)
	assert.NoError(t, err)

	expectedErr := errors.New("not enought arguments")

	result, err := db.Execute(ctx, "GET")
	assert.Nil(t, result)
	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestInvalidCommand(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := getMockedLogger()

	db, err := NewDatabase(compute, storage, logger)
	assert.NoError(t, err)

	expectedErr := errors.New("unknown command type")

	result, err := db.Execute(ctx, "gettt tttest")
	assert.Nil(t, result)
	assert.ErrorContains(t, err, expectedErr.Error())
}

func getMockedCompute(t *testing.T) *compute.Compute {
	compute, err := compute.NewCompute(getMockedLogger())
	assert.NoError(t, err)

	return compute
}

func getMockedLogger() *slog.Logger {
	buf := new(bytes.Buffer)
	return slog.New(slog.NewTextHandler(buf, nil))
}
