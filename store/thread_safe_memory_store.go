package store

import "sync"

type ThreadSafeMemoryStore[T any] struct {
	mu    *sync.RWMutex
	store map[string]T
}

func NewThreadSafeMemory[T any]() *ThreadSafeMemoryStore[T] {
	return &ThreadSafeMemoryStore[T]{
		mu:    &sync.RWMutex{},
		store: map[string]T{},
	}
}

func (ms *ThreadSafeMemoryStore[T]) Put(key string, value T) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.store[key] = value
}

func (ms *ThreadSafeMemoryStore[T]) Get(key string) T {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.store[key]
}

func (ms *ThreadSafeMemoryStore[T]) Delete(key string) bool {
	isDeleted := ms.Exists(key)
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.store, key)
	return isDeleted
}

func (ms *ThreadSafeMemoryStore[T]) Exists(key string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	_, ok := ms.store[key]
	return ok
}

func (ms *ThreadSafeMemoryStore[T]) GetAll() map[string]T {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	copyMap := make(map[string]T)
	for key, value := range ms.store {
		copyMap[key] = value
	}
	return copyMap
}
