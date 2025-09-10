package backend

import (
	"context"
	"sync"
)

type MemStore struct {
	mu sync.RWMutex
	m  map[string]item
}

type item struct {
	tokens     int64
	lastRefill int64
}

func NewMemStore() *MemStore {
	return &MemStore{m: make(map[string]item)}
}

func (s *MemStore) Get(_ context.Context, key string) (int64, int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.m[key]
	if !ok {
		return 0, 0, nil
	}
	return it.tokens, it.lastRefill, nil
}

func (s *MemStore) Set(_ context.Context, key string, tokens int64, lastRefill int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = item{tokens: tokens, lastRefill: lastRefill}
	return nil
}
