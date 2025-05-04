package wal

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"kdb/internal/utils"
)

func TestNewWritterEmptyWALDirectory(t *testing.T) {
	expectedErr := errEmptyWALDirectory

	_, err := newWritter("", nil)

	assert.ErrorIs(t, err, expectedErr)
}

func TestNewWritterInvalidLogger(t *testing.T) {
	expectedErr := errInvalidLogger

	_, err := newWritter("test", nil)

	assert.ErrorIs(t, err, expectedErr)
}

func TestNewWritterSuccess(t *testing.T) {
	_, err := newWritter("test", utils.NewMockedLogger())
	assert.NoError(t, err)
}

func TestBuildMessage(t *testing.T) {
	writter, err := newWritter("test", utils.NewMockedLogger())
	assert.NoError(t, err)

	rows := []string{
		"SET test test",
		"DEL test test",
		"SET test1 test1",
		"SET test2 test2",
		"SET test3 test3",
	}

	expected := "SET test test\nDEL test test\nSET test1 test1\nSET test2 test2\nSET test3 test3\n"

	actual := writter.buildMessage(rows)

	assert.Equal(t, expected, actual)
}
