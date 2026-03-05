package concurrency

import (
	"sync"

	"github.com/samber/lo"
)

type (
	SyncMap[T comparable, E any] struct {
		mu sync.Mutex
		m  map[T]E
	}
)

func NewSyncMap[T comparable, E any](m map[T]E) *SyncMap[T, E] {
	return &SyncMap[T, E]{
		m: m,
	}
}

func (m *SyncMap[T, E]) Get(idx T) (E, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, found := m.m[idx]
	return e, found
}

func (m *SyncMap[T, E]) Set(idx T, e E) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.m == nil {
		m.m = make(map[T]E)
	}
	m.m[idx] = e
}

func (m *SyncMap[T, E]) Delete(idx T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, idx)
}

func (m *SyncMap[T, E]) Keys() []T {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := lo.Keys(m.m)
	return keys
}
