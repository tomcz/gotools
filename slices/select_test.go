package slices

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	selector := func(v int) bool { return v%2 == 0 }

	src := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, []int{2, 4, 6}, Select(src, selector))

	src = []int{1, 3, 5, 7}
	assert.Len(t, Select(src, selector), 0)
}

func TestSelectErr(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	selector := func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err := SelectErr(src, selector)
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, actual)

	src = []int{1, 3, 5, 7}
	selector = func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err = SelectErr(src, selector)
	assert.NoError(t, err)
	assert.Len(t, actual, 0)

	src = []int{1, 2, 3, 4, 5, 6}
	selector = func(v int) (bool, error) { return false, errors.New("test error") }
	actual, err = SelectErr(src, selector)
	assert.Error(t, err, "test error")
	assert.Nil(t, actual)
}
