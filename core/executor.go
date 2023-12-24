package core

import (
	"errors"
	"strings"
)

func (c *Command) Execute() (any, error) {
	switch strings.ToUpper(c.Cmd) {
	case "PING":
		return executePing(c.Args...)
	case "ECHO":
		return executeEcho(c.Args...)
	default:
		return nil, errors.New("Err unsuported command")
	}
}

func executePing(args ...any) (any, error) {
	if len(args) > 1 {
		return nil, errors.New("ERR wrong number of arguments")
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return "PONG", nil
}

func executeEcho(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, errors.New("ERR wrong number of arguments")
	}
	return args[0], nil
}
