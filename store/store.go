package store

type Storer interface {
	Put(string, any)
	Get(string) any
	Del(string)
	Exists(string) bool
}
