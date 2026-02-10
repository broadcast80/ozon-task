package inmemory

import (
	"context"
	"sync"

	"github.com/broadcast80/ozon-task/internal/pkg/models"
)

type repository struct {
	aliasToURL map[string]string
	urlToAlias map[string]string
	mu         sync.RWMutex
}

func New(storeSize int) *repository {
	return &repository{
		aliasToURL: make(map[string]string, storeSize),
		urlToAlias: make(map[string]string, storeSize),
		mu:         sync.RWMutex{},
	}
}

func (r *repository) Create(ctx context.Context, url string, alias string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.aliasToURL[alias]; ok {
		return models.ErrDuplicate
	}

	r.aliasToURL[alias] = url
	r.urlToAlias[url] = alias

	return nil

}

func (r *repository) Get(ctx context.Context, alias string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.aliasToURL[alias]
	if !ok {
		return "", models.ErrNotFound
	}

	return url, nil
}

func (r *repository) URLExists(ctx context.Context, url string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.urlToAlias[url]
	return ok, nil
}
