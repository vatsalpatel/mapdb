package core

import (
	"log"
	"os"

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
	engine := &Engine{
		storage,
	}
	err := engine.load()
	if err != nil {
		log.Println(err)
	}
	return engine
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

func (e *Engine) load() error {
	data, err := os.ReadFile("dump.rdb")
	if err != nil {
		return err
	}

	var key, value, expiry []byte
	var current int

	for i := 0; i < len(data); i++ {
		switch {
		case data[i] == ',':
			current++
		case data[i] == '\n':
			_, err = e.execSet(string(key), string(value), string(expiry))
			if err != nil {
				log.Println(err)
				continue
			}
			key, value, expiry = []byte{}, []byte{}, []byte{}
			current = 0
		default:
			switch current {
			case 0:
				key = append(key, data[i])
			case 1:
				value = append(value, data[i])
			case 2:
				expiry = append(expiry, data[i])
			}
		}
	}
	return nil
}
