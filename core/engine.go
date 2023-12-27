package core

import (
	"github.com/vatsalpatel/radish/store"
)

type IEngine interface {
	store.Storer[*Item]
	Handle([]byte) []byte
}

type Engine struct {
	store.Storer[*Item]
}

type Item struct {
	value  any
	expiry int64
}

func NewEngine(storage store.Storer[*Item]) *Engine {
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
