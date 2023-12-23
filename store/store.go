package store

type Storer interface {
	Put(string, any)
	Get(string) any
	Delete(string)
	Exists(string) bool
}
