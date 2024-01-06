package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrWrongNumberOfArgs = errors.New("ERR wrong number of arguments")
	ErrWrongTypeOfArgs   = errors.New("ERR wrong type of arguments")
	ErrValueNotInteger   = errors.New("ERR value is not an integer or out of range")
)

func (e *Engine) execute(cmd *Command) (any, error) {
	cmd.Cmd = strings.ToUpper(cmd.Cmd)
	e.logCommand(cmd)

	switch cmd.Cmd {
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
	case "INCR":
		return e.execIncr(cmd.Args...)
	case "DECR":
		return e.execDecr(cmd.Args...)
	case "SAVE":
		return e.execSave(cmd.Args...)
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

func (e *Engine) logCommand(cmd *Command) {
	logLine := strings.Builder{}
	logLine.Write([]byte(time.Now().UTC().Format(time.RFC3339)))
	logLine.Write([]byte(" "))
	logLine.Write([]byte(cmd.Cmd))
	for _, arg := range cmd.Args {
		logLine.Write([]byte(" "))
		logLine.Write([]byte(fmt.Sprintf("%v", arg)))
	}
	err := e.PersistentStorer.Write([]byte(logLine.String()))
	if err != nil {
		log.Println("cannot write to log", err)
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
		if expiry != -1 {
			expiry = time.Now().UnixMilli() + expiry*1000
		}
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

func (e *Engine) execExpire(args ...any) (int64, error) {
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
	value, ok := args[1].(string)
	if !ok {
		return 0, ErrWrongTypeOfArgs
	}
	expiry, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, ErrWrongTypeOfArgs
	}
	item.expiry = time.Now().UnixMilli() + expiry*1000
	return 1, nil
}

func (e *Engine) execTTL(args ...any) (int64, error) {
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
	return int64(item.expiry-time.Now().UnixMilli()) / 1000, nil
}

func (e *Engine) execIncr(args ...any) (string, error) {
	if len(args) != 1 {
		return "", ErrWrongNumberOfArgs
	}

	key, ok := args[0].(string)
	if !ok {
		return "", ErrWrongTypeOfArgs
	}

	item, ok := e.getItem(key)
	if !ok {
		item = &Item{
			value:  "0",
			expiry: -1,
		}
	}

	valueStr, ok := item.value.(string)
	if !ok {
		return "", ErrValueNotInteger
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return "", ErrValueNotInteger
	}

	value++
	valueStr = strconv.FormatInt(value, 10)

	e.Storer.Put(key, &Item{
		value:  valueStr,
		expiry: item.expiry,
	})

	return valueStr, nil
}

func (e *Engine) execDecr(args ...any) (string, error) {
	if len(args) != 1 {
		return "", ErrWrongNumberOfArgs
	}

	key, ok := args[0].(string)
	if !ok {
		return "", ErrWrongTypeOfArgs
	}

	item, ok := e.getItem(key)
	if !ok {
		item = &Item{
			value:  "0",
			expiry: -1,
		}
	}

	valueStr, ok := item.value.(string)
	if !ok {
		return "", ErrValueNotInteger
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return "", ErrValueNotInteger
	}

	value--
	valueStr = strconv.FormatInt(value, 10)

	e.Storer.Put(key, &Item{
		value:  valueStr,
		expiry: item.expiry,
	})

	return valueStr, nil
}

func (e *Engine) execSave(args ...any) (string, error) {
	data := e.Storer.GetAll()

	// writes data to file called "dump.rdb" in format:
	// key1,value1,expiry1
	// key2,value2,expiry2
	// ...
	// keyN,valueN,expiryN

	bytes := []byte{}
	for key, value := range data {
		bytes = append(bytes, []byte(key)...)
		bytes = append(bytes, []byte(",")...)
		bytes = append(bytes, []byte(fmt.Sprintf("%v", value.value))...)
		bytes = append(bytes, []byte(",")...)
		bytes = append(bytes, []byte(fmt.Sprintf("%v", value.expiry))...)
		bytes = append(bytes, []byte("\n")...)
	}

	file, err := os.Create("dump.rdb")
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("Cannot close file", err)
		}
	}()
	if err != nil {
		return "", err
	}
	file.Write(bytes)

	return "OK", nil
}
