package sync

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllSettledWithoutErrors(t *testing.T) {
	f1 := NewFuture(func() (string, *error) {
		return "hello", nil
	})
	f2 := NewFuture(func() (string, *error) {
		return "world", nil
	})

	all := AllSettled(f1, f2)
	results, err := all.Wait()

	assert.Equal(t, []string{"hello", "world"}, results)
	assert.Nil(t, err)
}

func TestAllSettledWithErrors(t *testing.T) {
	f1 := NewFuture(func() (string, *error) {
		return "hello", nil
	})
	f2 := NewFuture(func() (string, *error) {
		err := errors.New("an Error")
		return "", &err
	})

	all := AllSettled(f1, f2)
	results, err := all.Wait()

	assert.Equal(t, []string{"hello", ""}, results)
	assert.Equal(t, "0: nil\n1: an Error\n", err.Error())
	assert.Nil(t, err.Errors[0])
	assert.Equal(t, "an Error", (*err.Errors[1]).Error())
}
