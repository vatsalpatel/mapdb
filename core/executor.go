package core

import (
	"errors"
	"strings"
)

var (
	ErrWrongNumberOfArgs = errors.New("ERR wrong number of arguments")
	ErrWrongTypeOfArgs   = errors.New("ERR wrong type of arguments")
)

func (e *Engine) execute(cmd *Command) (any, error) {
	switch strings.ToUpper(cmd.Cmd) {
	case "PING":
		return e.execPing(cmd.Args...)
	case "ECHO":
		return e.execEcho(cmd.Args...)
	case "SET":
		return e.execSet(cmd.Args...)
	case "GET":
		return e.execGet(cmd.Args...)
	case "DEL":
		return e.execDelete(cmd.Args...)
	case "EXISTS":
		return e.execExists(cmd.Args...)
	default:
		return nil, errors.New("Err unsuported command")
	}
}

func (e *Engine) execPing(args ...any) (any, error) {
	if len(args) > 1 {
		return nil, ErrWrongNumberOfArgs
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return "PONG", nil
}

func (e *Engine) execEcho(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs
	}
	return args[0], nil
}

func (e *Engine) execSet(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, ErrWrongNumberOfArgs
	}
	key, ok := args[0].(string)
	if !ok {
		return nil, ErrWrongTypeOfArgs
	}
	oldValue, err := e.execGet(key)
	if err != nil {
		return nil, err
	}
	e.Storer.Put(key, args[1])
	return oldValue, nil
}

func (e *Engine) execGet(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, ErrWrongNumberOfArgs
	}
	key, ok := args[0].(string)
	if !ok {
		return nil, ErrWrongTypeOfArgs
	}
	return e.Storer.Get(key), nil
}

func (e *Engine) execDelete(args ...any) (int, error) {
	if len(args) < 1 {
		return 0, ErrWrongNumberOfArgs
	}
	count := 0
	for _, arg := range args {
		key, ok := arg.(string)
		if !ok {
			return 0, ErrWrongTypeOfArgs
		}
		isDeleted := e.Storer.Delete(key)
		if isDeleted {
			count++
		}
	}
	return count, nil
}

func (e *Engine) execExists(args ...any) (int, error) {
	if len(args) < 1 {
		return 0, ErrWrongNumberOfArgs
	}
	count := 0
	for _, arg := range args {
		key, ok := arg.(string)
		if !ok {
			return 0, ErrWrongTypeOfArgs
		}
		if e.Storer.Exists(key) {
			count++
		}
	}
	return count, nil
}
