package store

type Storer interface {
	Put(string, any)
	Get(string) any
	Delete(string) bool
	Exists(string) bool
}
