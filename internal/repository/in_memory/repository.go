package inmemory

import (
	"context"
	"fmt"
	"sync"
)

type repository struct {
	store map[string]string
	mu    sync.RWMutex
}

func New(storeSize int) *repository {
	return &repository{
		store: make(map[string]string, storeSize),
		mu:    sync.RWMutex{},
	}
}

func (r *repository) Create(ctx context.Context, url string, alias string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[alias]; ok {
		return fmt.Errorf("url already exists")
	}

	// проверка существования такого алиаса

	r.store[alias] = url
	return nil

}

func (r *repository) Get(ctx context.Context, alias string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.store[alias]
	if !ok {
		return "", fmt.Errorf("not found")
	}

	return url, nil
}
