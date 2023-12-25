package store

import "sync"

type ThreadSafeMemoryStore struct {
	mu    *sync.RWMutex
	store *map[string]any
}

func NewThreadSafeMemory() *ThreadSafeMemoryStore {
	return &ThreadSafeMemoryStore{
		mu:    &sync.RWMutex{},
		store: &map[string]any{},
	}
}

func (ms *ThreadSafeMemoryStore) Put(key string, value any) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	(*ms.store)[key] = value
}

func (ms *ThreadSafeMemoryStore) Get(key string) any {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return (*ms.store)[key]
}

func (ms *ThreadSafeMemoryStore) Delete(key string) bool {
	isDeleted := ms.Exists(key)
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(*ms.store, key)
	return isDeleted
}

func (ms *ThreadSafeMemoryStore) Exists(key string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	_, ok := (*ms.store)[key]
	return ok
}
