package core

import (
	"github.com/vatsalpatel/radish/store"
)

type IEngine interface {
	store.Storer
	Execute([]byte) ([]byte, error)
}

type Engine struct {
	store.Storer
}

func NewEngine(store store.Storer) *Engine {
	return &Engine{
		store,
	}
}

func (e *Engine) Execute(input []byte) ([]byte, error) {
	return []byte("+OK\r\n"), nil
}
