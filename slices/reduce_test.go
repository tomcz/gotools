package slices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	assert.Equal(t, 10, Reduce([]int{1, 2, 3, 4}, 0, func(val int, acc int) int { return val + acc }))
}

func TestReduceErr(t *testing.T) {
	input := []int{1, 2, 3, 4}

	actual, err := ReduceErr(input, 0, func(val int, acc int) (int, error) { return val + acc, nil })
	if assert.NoError(t, err) {
		assert.Equal(t, 10, actual)
	}

	actual, err = ReduceErr(input, 0, func(val int, acc int) (int, error) { return 0, errors.New("test error") })
	if assert.EqualError(t, err, "test error") {
		assert.Equal(t, 0, actual)
	}
}
