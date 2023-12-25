package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vatsalpatel/radish/store"
)

func TestMemoryStore(t *testing.T) {
	ms := store.NewMemory()
	ms.Put("foo", "bar")
	assert.Equal(t, ms.Get("foo"), "bar")
	assert.Equal(t, ms.Exists("foo"), true)

	ms.Delete("foo")
	assert.Equal(t, ms.Get("foo"), nil)
	assert.Equal(t, ms.Exists("foo"), false)
}
