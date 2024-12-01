package sync

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIteratesOverSlice(t *testing.T) {
	future := ParallelMap([]int{1, 2, 3}, func(i int) (int, *error) {
		return i + 1, nil
	})

	result, _ := future.Wait()

	assert.Equal(t, []int{2, 3, 4}, result)
}

func TestReturnsError(t *testing.T) {
	future := ParallelMap([]int{1, 2, 3}, func(i int) (int, *error) {
		if i == 2 {
			err := errors.New("some error")
			return 0, &err
		} else {
			return i + 1, nil
		}
	})

	result, err := future.Wait()

	assert.Nil(t, result)
	assert.Equal(t, *err, errors.New("some error"))
}
