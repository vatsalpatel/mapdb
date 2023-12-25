package store

type MemoryStore struct {
	store *map[string]any
}

func NewMemory() *MemoryStore {
	return &MemoryStore{
		store: &map[string]any{},
	}
}

func (ms *MemoryStore) Put(key string, value any) {
	(*ms.store)[key] = value
}

func (ms *MemoryStore) Get(key string) any {
	return (*ms.store)[key]
}

func (ms *MemoryStore) Delete(key string) bool {
	isDeleted := ms.Exists(key)
	delete(*ms.store, key)
	return isDeleted
}

func (ms *MemoryStore) Exists(key string) bool {
	_, ok := (*ms.store)[key]
	return ok
}
