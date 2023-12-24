package core

import (
	"errors"
)

func Deserialize(input []byte) (any, int, error) {
	if len(input) == 0 {
		return nil, 0, errors.New("input is empty")
	}

	switch input[0] {
	case '+':
		return readSimpleString(input)
	case '-':
		return readError(input)
	case ':':
		return readInteger(input)
	case '$':
		return readBulkString(input)
	case '*':
		return readArray(input)
	default:
		return nil, 0, errors.New("type is invalid")
	}
}

func readLength(input []byte) (int, int) {
	pos, length := 1, 0
	for ; input[pos] != '\r'; pos++ {
		length = length*10 + int(input[pos]-'0')
	}

	return length, pos + 2
}

func readSimpleString(input []byte) ([]byte, int, error) {
	pos := 1
	str := make([]byte, 0)
	for ; input[pos] != '\r'; pos++ {
		str = append(str, input[pos])
	}
	return str, pos + 2, nil
}

func readBulkString(input []byte) ([]byte, int, error) {
	length, delta := readLength(input)
	if length < 0 {
		return nil, delta, nil
	}

	str := make([]byte, length)
	copy(str, input[delta:delta+length])
	return str, delta + length + 2, nil
}

func readInteger(input []byte) (int, int, error) {
	pos, num, sign := 1, 0, 1
	if input[pos] == '+' {
		pos++
	}
	if input[pos] == '-' {
		sign = -1
		pos++
	}
	for ; input[pos] != '\r'; pos++ {
		num = num*10 + int(input[pos]-'0')
	}
	return sign * num, pos + 2, nil
}

func readArray(input []byte) ([]any, int, error) {
	length, offset := readLength(input)
	if length <= 0 {
		return []any{}, offset, nil
	}

	arr := make([]any, length)
	for i := 0; i < length; i++ {
		data, delta, err := Deserialize(input[offset:])
		if err != nil {
			return nil, delta, err
		}
		arr[i] = data
		offset += delta
	}

	return arr, 0, nil
}

func readError(input []byte) ([]byte, int, error) {
	return readSimpleString(input)
}

func serialize(any any) ([]byte, error) {
	return nil, nil
}
