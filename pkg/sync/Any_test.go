package sync

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReturnsFirstSuccesfulFuture(t *testing.T) {
	f1 := NewFuture(func() (string, *error) {
		time.Sleep(time.Second)
		return "hello", nil
	})
	f2 := NewFuture(func() (string, *error) {
		return "world", nil
	})

	f := Any(f1, f2)
	results, err := f.Wait()

	assert.Equal(t, "world", results)
	assert.Nil(t, err)
}

func TestReturnsAllErrors(t *testing.T) {
	f1 := NewFuture(func() (string, *error) {
		err := errors.New("error1")
		return "ignored", &err
	})
	f2 := NewFuture(func() (string, *error) {
		err := errors.New("error2")
		return "also ignored", &err

	})

	f := Any(f1, f2)
	results, err := f.Wait()

	assert.Equal(t, "", results)
	assert.Contains(t, "0: error1\n1: error2\n", err.Error())
}
