package engine

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	ctx := context.Background()
	engine := NewEngine()

	data := getTestData()

	// test concurrent writing
	for i := range data {
		go func(item string) {
			err := engine.Set(ctx, item, item)
			assert.Nil(t, err)
		}(data[i])
	}

	// test getting
	for i := range data {
		go func(item string) {
			i, err := engine.Get(ctx, item)
			assert.Nil(t, err)
			assert.Equal(t, item, i)
		}(data[i])
	}

	// test deliting
	wg := &sync.WaitGroup{}
	for i := range data {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()

			err := engine.Del(ctx, item)
			assert.Nil(t, err)
		}(data[i])
	}
	wg.Wait()

	assert.Equal(t, 0, len(engine.m))
}

func getTestData() []string {
	sl := make([]string, 0, 1000)

	for i := range 1000 {
		sl = append(sl, strconv.Itoa(i))
	}

	return sl
}
