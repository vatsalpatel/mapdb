package store

type MemoryStore[T any] struct {
	store map[string]T
}

func NewMemory[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{
		store: map[string]T{},
	}
}

func (ms *MemoryStore[T]) Put(key string, value T) {
	ms.store[key] = value
}

func (ms *MemoryStore[T]) Get(key string) T {
	return ms.store[key]
}

func (ms *MemoryStore[T]) Delete(key string) bool {
	isDeleted := ms.Exists(key)
	delete(ms.store, key)
	return isDeleted
}

func (ms *MemoryStore[T]) Exists(key string) bool {
	_, ok := ms.store[key]
	return ok
}

func (ms *MemoryStore[T]) GetAll() map[string]T {
	copyMap := make(map[string]T)
	for key, value := range ms.store {
		copyMap[key] = value
	}
	return copyMap
}
