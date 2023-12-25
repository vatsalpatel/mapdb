package core

import (
	"github.com/vatsalpatel/radish/store"
)

type IEngine interface {
	store.Storer
	Handle([]byte) []byte
}

type Engine struct {
	store.Storer
}

func NewEngine(storage store.Storer) *Engine {
	return &Engine{
		storage,
	}
}

func (e *Engine) Handle(input []byte) []byte {
	deserialized, err := DeserializeArray(input)
	if err != nil {
		return Serialize(err)
	}
	command := &Command{
		Cmd:  deserialized[0].(string),
		Args: deserialized[1:],
	}
	result, err := e.execute(command)
	if err != nil {
		return Serialize(err)
	}
	return Serialize(result)
}
