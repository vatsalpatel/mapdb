package store

type Storer[T any] interface {
	Put(string, T)
	Get(string) T
	Delete(string) bool
	Exists(string) bool
	GetAll() map[string]T
}

type PersistentStorer interface {
	Write([]byte) error
	ReadAll() ([]byte, error)
	Clear() error
}
