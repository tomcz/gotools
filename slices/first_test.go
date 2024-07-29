package slices

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestFirst(t *testing.T) {
	selector := func(v int) bool { return v%2 == 0 }

	src := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, 2, First(src, selector))

	src = []int{1, 3, 5, 7}
	assert.Equal(t, 0, First(src, selector))
}

func TestFirstErr(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	selector := func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err := FirstErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, 2, actual)

	src = []int{1, 3, 5, 7}
	selector = func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err = FirstErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, 0, actual)

	src = []int{1, 2, 3, 4, 5, 6}
	selector = func(v int) (bool, error) { return false, errors.New("test error") }
	actual, err = FirstErr(src, selector)
	assert.Error(t, err, "test error")
	assert.Equal(t, 0, actual)
}
