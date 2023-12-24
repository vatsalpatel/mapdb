package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vatsalpatel/radish/core"
)

func TestSimleString(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte("+OK\r\n"), []byte("OK")},
		{[]byte("+PONG\r\n"), []byte("PONG")},
		{[]byte("+\r\n"), []byte("")},
	}
	for tc := range testCases {
		actual, _, _ := core.Deserialize(testCases[tc].input)
		assert.Equal(t, testCases[tc].expected, actual)
	}
}

func TestBulkString(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte("$6\r\nfoobar\r\n"), []byte("foobar")},
		{[]byte("$0\r\n\r\n"), []byte{}},
	}
	for tc := range testCases {
		actual, _, _ := core.Deserialize(testCases[tc].input)
		assert.Equal(t, testCases[tc].expected, actual)
	}
}

func TestArray(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    []byte
		expected []interface{}
	}{
		{[]byte("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"), []any{[]byte("foo"), []byte("bar")}},
		{[]byte("*0\r\n"), []any{}},
		{[]byte("*4\r\n:1\r\n:2\r\n+qwe\r\n$6\r\nfoobar\r\n"), []any{1, 2, []byte("qwe"), []byte("foobar")}},
	}
	for tc := range testCases {
		actual, _, err := core.Deserialize(testCases[tc].input)
		assert.Nil(t, err)
		assert.Equal(t, testCases[tc].expected, actual)
	}
}

func TestInteger(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    []byte
		expected int
	}{
		{[]byte(":0\r\n"), 0},
		{[]byte(":-12\r\n"), -12},
		{[]byte(":123\r\n"), 123},
		{[]byte(":+123\r\n"), 123},
	}
	for tc := range testCases {
		actual, _, _ := core.Deserialize(testCases[tc].input)
		assert.Equal(t, testCases[tc].expected, actual)
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte("-Error message\r\n"), []byte("Error message")},
		{[]byte("-\r\n"), []byte{}},
	}
	for tc := range testCases {
		actual, _, _ := core.Deserialize(testCases[tc].input)
		assert.Equal(t, testCases[tc].expected, actual)
	}
}
