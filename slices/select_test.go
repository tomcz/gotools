package slices

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestSelect(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	selector := func(v int) bool { return v%2 == 0 }
	assert.DeepEqual(t, []int{2, 4, 6}, Select(src, selector))
}

func TestSelectErr(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	selector := func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err := SelectErr(src, selector)
	assert.NilError(t, err)
	assert.DeepEqual(t, []int{2, 4, 6}, actual)

	selector = func(v int) (bool, error) { return false, errors.New("test error") }
	actual, err = SelectErr(src, selector)
	assert.Error(t, err, "test error")
	assert.Assert(t, is.Nil(actual))
}
