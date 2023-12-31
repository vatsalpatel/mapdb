package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vatsalpatel/radish/store"
)

func TestThreadSafeMemoryStore(t *testing.T) {
	ms := store.NewThreadSafeMemory[any]()
	ms.Put("foo", "bar")
	assert.Equal(t, ms.Get("foo"), "bar")
	assert.Equal(t, ms.Exists("foo"), true)

	ms.Delete("foo")
	assert.Equal(t, ms.Get("foo"), nil)
	assert.Equal(t, ms.Exists("foo"), false)

	ms.Put("foo", "bar")
	ms.Put("bar", "baz")
	assert.Equal(t, ms.GetAll(), map[string]any{"foo": "bar", "bar": "baz"})
}
