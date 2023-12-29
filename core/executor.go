package core

import (
	"errors"
	"strconv"
	"strings"
	"time"
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
	case "EXPIRE":
		return e.execExpire(cmd.Args...)
	case "TTL":
		return e.execTTL(cmd.Args...)
	default:
		return nil, errors.New("Err unsuported command")
	}
}

func (e *Engine) getItem(key string) (*Item, bool) {
	item := e.Storer.Get(key)
	if item == nil {
		return nil, false
	}
	if item.expiry > 0 && time.Now().UTC().UnixMilli() > item.expiry {
		e.Storer.Delete(key)
		return nil, false
	}
	return item, true
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
	if len(args) != 2 && len(args) != 3 {
		return nil, ErrWrongNumberOfArgs
	}

	key, ok := args[0].(string)
	if !ok {
		return nil, ErrWrongTypeOfArgs
	}
	var expiry int64 = -1
	if len(args) == 3 {
		var err error
		expiry, err = strconv.ParseInt(args[2].(string), 10, 64)
		if err != nil {
			return nil, ErrWrongTypeOfArgs
		}
		expiry = time.Now().UnixMilli() + expiry*1000
	}

	oldItem, exists := e.getItem(key)
	var oldValue any = "<nil>"
	if exists {
		oldValue = oldItem.value
	}

	e.Storer.Put(key, &Item{
		value:  args[1],
		expiry: expiry,
	})

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
	item, exists := e.getItem(key)
	if exists == false {
		return "<nil>", nil
	}
	return item.value, nil
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
		if _, exists := e.getItem(key); exists {
			count++
		}
	}
	return count, nil
}

func (e *Engine) execExpire(args ...any) (int, error) {
	if len(args) != 2 {
		return 0, ErrWrongNumberOfArgs
	}
	key, ok := args[0].(string)
	if !ok {
		return 0, ErrWrongTypeOfArgs
	}
	item, exists := e.getItem(key)
	if exists == false {
		return 0, nil
	}
	expiry, err := strconv.ParseInt(args[1].(string), 10, 64)
	if err != nil {
		return 0, ErrWrongTypeOfArgs
	}
	item.expiry = time.Now().UnixMilli() + expiry*1000
	return 1, nil
}

func (e *Engine) execTTL(args ...any) (int, error) {
	if len(args) != 1 {
		return 0, ErrWrongNumberOfArgs
	}
	key, ok := args[0].(string)
	if !ok {
		return 0, ErrWrongTypeOfArgs
	}
	item, exists := e.getItem(key)
	if exists == false {
		return -2, nil
	}
	if item.expiry == -1 {
		return -1, nil
	}
	return int(item.expiry-time.Now().UnixMilli()) / 1000, nil
}
