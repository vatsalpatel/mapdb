package store

type Storer[T any] interface {
	Put(string, T)
	Get(string) T
	Delete(string) bool
	Exists(string) bool
}
