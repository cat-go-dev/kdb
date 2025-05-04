package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"kdb/internal/utils"
	"kdb/mocks"
)

func TestNewStorageWithEmptyLogger(t *testing.T) {
	engine := mocks.NewEngineLayer(t)
	client, err := NewStorage(engine, nil)

	assert.Nil(t, client)
	assert.ErrorIs(t, err, errInvalidLogger)
}

func TestGetSuccess(t *testing.T) {
	ctx := context.Background()
	engine := mocks.NewEngineLayer(t)

	key := "val"
	value := "true"

	logger := utils.NewMockedLogger()

	st, _ := NewStorage(engine, logger)

	engine.EXPECT().Get(ctx, key).Return(value, nil)

	actual, err := st.Get(ctx, key)

	assert.Nil(t, err)
	assert.Equal(t, value, actual)
}

func TestGetError(t *testing.T) {
	ctx := context.Background()
	engine := mocks.NewEngineLayer(t)

	key := "val"
	value := ""
	expectedErr := fmt.Errorf("engine error")

	st, _ := NewStorage(engine, utils.NewMockedLogger())

	engine.EXPECT().Get(ctx, key).Return(value, expectedErr)

	actual, err := st.Get(ctx, key)

	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, value, actual)
}

func TestSetSuccess(t *testing.T) {
	ctx := context.Background()
	engine := mocks.NewEngineLayer(t)

	key := "val"
	value := "true"

	st, _ := NewStorage(engine, utils.NewMockedLogger())

	engine.EXPECT().Set(ctx, key, value).Return(nil)

	err := st.Set(ctx, key, value)

	assert.Nil(t, err)
}

func TestSetError(t *testing.T) {
	ctx := context.Background()
	engine := mocks.NewEngineLayer(t)

	key := "val"
	value := "true"
	expextedErr := fmt.Errorf("engine error")

	st, _ := NewStorage(engine, utils.NewMockedLogger())

	engine.EXPECT().Set(ctx, key, value).Return(expextedErr)

	err := st.Set(ctx, key, value)

	assert.ErrorIs(t, err, expextedErr)
}

func TestDelSuccess(t *testing.T) {
	ctx := context.Background()
	engine := mocks.NewEngineLayer(t)

	key := "val"

	st, _ := NewStorage(engine, utils.NewMockedLogger())

	engine.EXPECT().Del(ctx, key).Return(nil)

	err := st.Del(ctx, key)

	assert.Nil(t, err)
}

func TestDel(t *testing.T) {
	ctx := context.Background()
	engine := mocks.NewEngineLayer(t)

	key := "val"
	expextedErr := fmt.Errorf("engine error")

	st, _ := NewStorage(engine, utils.NewMockedLogger())

	engine.EXPECT().Del(ctx, key).Return(expextedErr)

	err := st.Del(ctx, key)

	assert.ErrorIs(t, err, expextedErr)
}
