package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vatsalpatel/mapdb/store"
)

func getMockPersistentStore() store.PersistentStorer {
	return &MockPersistentStore{}
}

func getMockEngine() *Engine {
	memoryStore := store.NewThreadSafeMemory[*Item]()
	persistentStore := getMockPersistentStore()
	return NewEngine(memoryStore, persistentStore)
}

func calcExpiey(expiry int64) int64 {
	if expiry == 0 {
		return -1
	}

	return time.Now().UTC().Unix() + expiry
}

func TestPing(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected string
		isError  bool
	}{
		{[]any{}, "PONG", false},
		{[]any{"Hello"}, "Hello", false},
		{[]any{"Hello", "World"}, "", true},
	}

	for _, tc := range testCases {
		result, err := e.execPing(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestEcho(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected string
		isError  bool
	}{
		{[]any{}, "", true},
		{[]any{"Hello"}, "Hello", false},
		{[]any{"Hello", "World"}, "", true},
	}

	for _, tc := range testCases {
		result, err := e.execEcho(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestSet(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected string
		isError  bool
	}{
		{[]any{}, "", true},
		{[]any{"Hello"}, "", true},
		{[]any{"Hello", "World"}, "<nil>", false},
		{[]any{"Hello", "World", "1000"}, "World", false},
		{[]any{"Hello", "World", "1000", "Extra"}, "", true},
	}

	for _, tc := range testCases {
		result, err := e.execSet(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected string
		isError  bool
	}{
		{[]any{}, "", true},
		{[]any{"f"}, "<nil>", false},
		{[]any{"a"}, "b", false},
	}

	e.execSet("a", "b")
	for _, tc := range testCases {
		result, err := e.execGet(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected int
		isError  bool
	}{
		{[]any{}, 0, true},
		{[]any{"f"}, 0, false},
		{[]any{"a"}, 1, false},
		{[]any{"a", "c", "f"}, 1, false},
	}

	e.execSet("a", "b")
	e.execSet("c", "d")
	for _, tc := range testCases {
		result, err := e.execDelete(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestExists(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected int
		isError  bool
	}{
		{[]any{}, 0, true},
		{[]any{"f"}, 0, false},
		{[]any{"a"}, 1, false},
		{[]any{"a", "c", "f"}, 2, false},
	}

	e.execSet("a", "b")
	e.execSet("c", "d")
	for _, tc := range testCases {
		result, err := e.execExists(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestExpire(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected int64
		isError  bool
	}{
		{[]any{}, 0, true},
		{[]any{"f"}, 0, true},
		{[]any{"a", "10"}, 10, false},
		{[]any{"c", "100"}, 100, false},
	}

	e.execSet("a", "b", "10")
	e.execSet("c", "d")
	for _, tc := range testCases {
		result, err := e.execExpire(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.NotEqual(t, 0, result)
			assert.LessOrEqual(t, (e.Storer.Get("a")).expiry/1000, calcExpiey(tc.expected)/1000)
		}
	}
}

func TestTTL(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected int64
		isError  bool
	}{
		{[]any{}, 0, true},
		{[]any{"f"}, -2, false},
		{[]any{"a"}, 10, false},
		{[]any{"c"}, -1, false},
	}

	e.execSet("a", "b", "10")
	e.execSet("c", "d")
	for _, tc := range testCases {
		result, err := e.execTTL(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.LessOrEqual(t, result, tc.expected)
		}
	}
}

func TestIncr(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected string
		isError  bool
	}{
		{[]any{}, "0", true},
		{[]any{"c"}, "0", true},
		{[]any{"a"}, "11", false},
		{[]any{"f"}, "1", false},
	}

	e.execSet("a", "10")
	e.execSet("c", "d")
	for _, tc := range testCases {
		result, err := e.execIncr(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

func TestDecr(t *testing.T) {
	t.Parallel()
	e := getMockEngine()
	testCases := []struct {
		args     []any
		expected string
		isError  bool
	}{
		{[]any{}, "0", true},
		{[]any{"c"}, "0", true},
		{[]any{"a"}, "9", false},
		{[]any{"f"}, "-1", false},
	}

	e.execSet("a", "10")
	e.execSet("c", "d")
	for _, tc := range testCases {
		result, err := e.execDecr(tc.args...)
		if tc.isError {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
}

type MockPersistentStore struct{}

func (m *MockPersistentStore) Write([]byte) error {
	return nil
}

func (m *MockPersistentStore) ReadAll() ([]byte, error) {
	return []byte{}, nil
}

func (m *MockPersistentStore) Clear() error {
	return nil
}
