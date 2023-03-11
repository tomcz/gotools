package slices

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestReduce(t *testing.T) {
	input := []int{1, 2, 3, 4}
	reducer := func(val int, acc int) int { return val + acc }
	assert.Equal(t, 10, Reduce(input, 0, reducer))
}

func TestReduceErr(t *testing.T) {
	input := []int{1, 2, 3, 4}

	reducer := func(val int, acc int) (int, error) { return val + acc, nil }
	actual, err := ReduceErr(input, 0, reducer)
	assert.NilError(t, err)
	assert.Equal(t, 10, actual)

	reducerErr := func(val int, acc int) (int, error) { return -1, errors.New("test error") }
	actual, err = ReduceErr(input, 0, reducerErr)
	assert.Error(t, err, "test error")
	assert.Equal(t, -1, actual)
}
