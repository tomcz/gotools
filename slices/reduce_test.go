package slices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if assert.NoError(t, err) {
		assert.Equal(t, 10, actual)
	}

	reducerErr := func(val int, acc int) (int, error) { return 0, errors.New("test error") }
	actual, err = ReduceErr(input, 0, reducerErr)
	if assert.EqualError(t, err, "test error") {
		assert.Equal(t, 0, actual)
	}
}
