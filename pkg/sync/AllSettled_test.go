package sync

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	data := []string{"value1", "value2", "value3"}
	futures := make([]*Future[int, error], len(data))
	for i, value := range data {
		f := NewFuture(func() (int, *error) {
			fmt.Println("Calling db with value", value)
			// Do some actual database calls here
			return i, nil
		})
		futures = append(futures, f)
	}

	results, errors := AllSettled(futures...).Wait()

	// Check the results and errors
}
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
