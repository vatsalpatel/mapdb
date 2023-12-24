package core

import (
	"errors"
	"fmt"
	"strings"
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

func readSimpleString(input []byte) (string, int, error) {
	pos := 1
	str := make([]byte, 0)
	for ; input[pos] != '\r'; pos++ {
		str = append(str, input[pos])
	}
	return string(str), pos + 2, nil
}

func readBulkString(input []byte) (string, int, error) {
	length, delta := readLength(input)
	if length < 0 {
		return "", delta, nil
	}

	str := make([]byte, length)
	copy(str, input[delta:delta+length])
	return string(str), delta + length + 2, nil
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

func readError(input []byte) (string, int, error) {
	return readSimpleString(input)
}

func Serialize(input any) ([]byte, error) {
	var builder strings.Builder
	switch input.(type) {
	case []byte:
		data := input.([]byte)
		builder.WriteString("$" + fmt.Sprintf("%v\r\n", len(data)) + string(data) + "\r\n")
	case string:
		builder.WriteString("+" + fmt.Sprintf("%v\r\n", input))
	case int:
		builder.WriteString(":" + fmt.Sprintf("%v\r\n", input))
	case error:
		builder.WriteString("-" + fmt.Sprintf("%v\r\n", input))
	case []any:
		data := input.([]any)
		builder.WriteString("*" + fmt.Sprintf("%v\r\n", len(data)))
		for _, item := range data {
			serialized, err := Serialize(item)
			if err != nil {
				return nil, err
			}
			builder.Write(serialized)
		}
	}
	return []byte(builder.String()), nil
}
