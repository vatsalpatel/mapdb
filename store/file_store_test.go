package store_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vatsalpatel/radish/store"
)

func TestFileStore(t *testing.T) {
	fs := store.NewFileStore("testfile.txt")

	// Test Append
	data := []byte("test data")
	err := fs.Write(data)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Test ReadAll
	result, err := fs.ReadAll()
	expected := append(data, '\r', '\n')
	assert.NoError(t, err)
	assert.Equal(t, result, expected)

	// Test Clear
	err = fs.Clear()
	assert.NoError(t, err)

	// Check that the file is empty after Clear
	result, err = fs.ReadAll()
	assert.NoError(t, err)
	assert.Len(t, result, 0)

	// Clean up the test file
	os.Remove("testfile.txt")
	assert.NoError(t, err)
}
