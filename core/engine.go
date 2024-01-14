package core

import (
	"log"

	"github.com/vatsalpatel/mapdb/store"
)

type IEngine interface {
	store.Storer[*Item]
	store.PersistentStorer
	Handle([]byte) []byte
	Shutdown() error
}

type Engine struct {
	store.Storer[*Item]
	store.PersistentStorer
}

type Item struct {
	value  any
	expiry int64
}

func NewEngine(memoryStorage store.Storer[*Item], logStorage store.PersistentStorer) *Engine {
	engine := &Engine{
		memoryStorage,
		logStorage,
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
		log.Println(command, err)
		return Serialize(err)
	}
	return Serialize(result)
}

func (e *Engine) load() error {
	data, err := e.PersistentStorer.ReadAll()
	if err != nil {
		return err
	}

	var key, value, expiry []byte
	var current int

	for i := 0; i < len(data); i++ {
		switch {
		case data[i] == ',':
			current++
		case data[i] == '\r':
			_, err = e.execSet(string(key), string(value), string(expiry))
			if err != nil {
				log.Println(err)
				continue
			}
			key, value, expiry = []byte{}, []byte{}, []byte{}
			current = 0
		case data[i] == '\n':
			continue
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

func (e *Engine) Shutdown() error {
	if _, err := e.execSave(); err != nil {
		return err
	}
	return nil
}
