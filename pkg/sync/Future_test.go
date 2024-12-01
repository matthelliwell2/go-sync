package sync

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWaitsForResult(t *testing.T) {
	f := NewFuture[string, error](func() (string, *error) {
		return "Hello world", nil
	})

	result, err := f.Wait()

	assert.Equal(t, "Hello world", result)
	assert.Nil(t, err)
}

func TestWaitsForError(t *testing.T) {
	f := NewFuture[string, error](func() (string, *error) {
		err := errors.New("error text")
		return "", &err
	})

	result, err := f.Wait()

	assert.Equal(t, "", result)
	assert.Error(t, errors.New("error text"), err)
}
