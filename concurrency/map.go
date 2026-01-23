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
	defer func() {
		m.mu.Unlock()
	}()
	e, found := m.m[idx]
	return e, found
}

func (m *SyncMap[T, E]) Set(idx T, e E) {
	m.mu.Lock()
	defer func() {
		m.mu.Unlock()
	}()
	m.m[idx] = e
}

func (m *SyncMap[T, E]) Delete(idx T) {
	m.mu.Lock()
	defer func() {
		m.mu.Unlock()
	}()
	delete(m.m, idx)
}

func (m *SyncMap[T, E]) Keys() []T {
	return lo.Keys(m.m)
}
