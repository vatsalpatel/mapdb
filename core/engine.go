package core

import (
	"github.com/vatsalpatel/radish/store"
)

type IEngine interface {
	store.Storer
	Execute([]byte) []byte
}

type Engine struct {
	store.Storer
}

func NewEngine(store store.Storer) *Engine {
	return &Engine{
		store,
	}
}

func (e *Engine) Execute(input []byte) []byte {
	deserialized, err := DeserializeArray(input)
	if err != nil {
		return Serialize(err)
	}
	command := &Command{
		Cmd:  deserialized[0].(string),
		Args: deserialized[1:],
	}
	result, err := command.Execute()
	if err != nil {
		return Serialize(err)
	}
	return Serialize(result)
}
