package database

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"kdb/internal/database/compute"
	"kdb/internal/utils"
	"kdb/mocks"
)

func TestNewDatabaseEmptyCompute(t *testing.T) {
	expectedErr := errInvalidCompute

	_, err := NewDatabase(nil, nil, nil, nil)

	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestNewDatabaseEmptyStorage(t *testing.T) {
	compute := getMockedCompute(t)

	expectedErr := errInvalidStorage

	_, err := NewDatabase(compute, nil, nil, nil)

	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestNewDatabaseEmptyLogger(t *testing.T) {
	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)

	expectedErr := errInvalidLogger

	_, err := NewDatabase(compute, storage, nil, nil)

	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestNotEnoughArguments(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := utils.NewMockedLogger()
	wal := mocks.NewWALLayer(t)

	db, err := NewDatabase(compute, storage, wal, logger)
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
	logger := utils.NewMockedLogger()
	wal := mocks.NewWALLayer(t)

	db, err := NewDatabase(compute, storage, wal, logger)
	assert.NoError(t, err)

	expectedErr := errors.New("unknown command type")

	result, err := db.Execute(ctx, "gettt tttest")
	assert.Nil(t, result)
	assert.ErrorContains(t, err, expectedErr.Error())
}

func TestWALShouldBeOnlyInSetAndGetMethods(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := utils.NewMockedLogger()
	wal := mocks.NewWALLayer(t)

	db, err := NewDatabase(compute, storage, wal, logger)
	assert.NoError(t, err)

	key := "test"
	val := "test"
	command := fmt.Sprintf("SET %s %s", key, val)

	storage.EXPECT().Set(ctx, key, "test").Return(nil)
	wal.EXPECT().Write(ctx, command).Return()

	_, err = db.Execute(ctx, command)
	assert.NoError(t, err)

	command = fmt.Sprintf("DEL %s", key)

	storage.EXPECT().Del(ctx, key).Return(nil)
	wal.EXPECT().Write(ctx, command).Return()

	_, err = db.Execute(ctx, command)
	assert.NoError(t, err)

	command = fmt.Sprintf("GET %s", key)

	storage.EXPECT().Get(ctx, key).Return(val, nil)

	_, err = db.Execute(ctx, command)
	assert.NoError(t, err)
}

func TestForGetCommandShouldBeCalledStorageGet(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := utils.NewMockedLogger()
	wal := mocks.NewWALLayer(t)

	db, err := NewDatabase(compute, storage, wal, logger)
	assert.NoError(t, err)

	key := "test"
	command := fmt.Sprintf("GET %s", key)

	storage.EXPECT().Get(ctx, key).Return("test", nil)

	_, err = db.Execute(ctx, command)
	assert.NoError(t, err)
}

func TestForSetCommandShouldBeCalledStorageGet(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := utils.NewMockedLogger()
	wal := mocks.NewWALLayer(t)

	db, err := NewDatabase(compute, storage, wal, logger)
	assert.NoError(t, err)

	key := "test"
	val := "test"
	command := fmt.Sprintf("SET %s %s", key, val)

	storage.EXPECT().Set(ctx, key, val).Return(nil)
	wal.EXPECT().Write(ctx, command).Return()

	_, err = db.Execute(ctx, command)
	assert.NoError(t, err)
}

func TestForDelCommandShouldBeCalledStorageGet(t *testing.T) {
	ctx := context.Background()

	compute := getMockedCompute(t)
	storage := mocks.NewStorageLayer(t)
	logger := utils.NewMockedLogger()
	wal := mocks.NewWALLayer(t)

	db, err := NewDatabase(compute, storage, wal, logger)
	assert.NoError(t, err)

	key := "test"
	command := fmt.Sprintf("DEL %s", key)

	storage.EXPECT().Del(ctx, key).Return(nil)
	wal.EXPECT().Write(ctx, command).Return()

	_, err = db.Execute(ctx, command)
	assert.NoError(t, err)
}

func getMockedCompute(t *testing.T) *compute.Compute {
	compute, err := compute.NewCompute(utils.NewMockedLogger())
	assert.NoError(t, err)

	return compute
}
