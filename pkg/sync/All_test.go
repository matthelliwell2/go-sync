package sync

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWaitsForAllSuccesses(t *testing.T) {
	f1 := NewFuture(func() (string, *error) {
		return "hello", nil
	})
	f2 := NewFuture(func() (string, *error) {
		return "world", nil
	})

	all := All(f1, f2)
	results, err := all.Wait()

	assert.Equal(t, []string{"hello", "world"}, results)
	assert.Nil(t, err)
}

func TestWaitsForFirstFailure(t *testing.T) {
	f1 := NewFuture(func() (string, *error) {
		// We shouldn't wait for this to complete as the second future fails straight away
		time.Sleep(time.Second * 100)
		return "hello", nil
	})
	f2 := NewFuture(func() (string, *error) {
		err := errors.New("an Error")
		return "world", &err
	})
	f3 := NewFuture(func() (string, *error) {
		return "world", nil
	})

	start := time.Now()
	all := All(f1, f2, f3)
	results, err := all.Wait()
	end := time.Now()

	assert.Nil(t, results)
	assert.Equal(t, "an Error", (*err).Error())
	assert.WithinDuration(t, start, end, time.Second)
}
