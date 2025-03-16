package engine

import (
	"context"
	"sync"
)

type Engine struct {
	mu *sync.Mutex
	m  map[string]string
}

const defaultMapSize = 1000

func NewEngine() *Engine {
	return &Engine{
		mu: &sync.Mutex{},
		m:  make(map[string]string, defaultMapSize),
	}
}

func (e *Engine) Get(ctx context.Context, key string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	value, exists := e.m[key]
	if !exists {
		return "", nil
	}

	return value, nil
}

func (e *Engine) Set(ctx context.Context, key, value string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.m[key] = value

	return nil
}

func (e *Engine) Del(ctx context.Context, key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.m, key)

	return nil
}
