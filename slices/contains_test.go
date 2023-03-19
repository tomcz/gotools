package slices

import (
	"errors"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestContains(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, true, Contains(src, 3))
	assert.Equal(t, false, Contains(src, 9))
}

func TestContainsAny(t *testing.T) {
	selector := func(v int) bool { return v%2 == 0 }

	src := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, true, ContainsAny(src, selector))

	src = []int{1, 3, 5, 7}
	assert.Equal(t, false, ContainsAny(src, selector))
}

func TestContainsAnyErr(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	selector := func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err := ContainsAnyErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, true, actual)

	src = []int{1, 3, 5, 7}
	selector = func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err = ContainsAnyErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, false, actual)

	src = []int{1, 2, 3, 4, 5, 6}
	selector = func(v int) (bool, error) { return false, errors.New("test error") }
	actual, err = ContainsAnyErr(src, selector)
	assert.Error(t, err, "test error")
	assert.Equal(t, false, actual)
}

func TestContainsAll(t *testing.T) {
	selector := func(s string) bool { return strings.HasPrefix(s, "test_") }

	src := []string{"test_one", "test_two", "test_three"}
	assert.Equal(t, true, ContainsAll(src, selector))

	src = []string{"test_one", "nope", "test_three"}
	assert.Equal(t, false, ContainsAll(src, selector))
}

func TestContainsAllErr(t *testing.T) {
	selector := func(s string) (bool, error) { return strings.HasPrefix(s, "test_"), nil }

	src := []string{"test_one", "test_two", "test_three"}
	actual, err := ContainsAllErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, true, actual)

	src = []string{"test_one", "nope", "test_three"}
	actual, err = ContainsAllErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, false, actual)

	selector = func(s string) (bool, error) { return false, errors.New("test error") }
	actual, err = ContainsAllErr(src, selector)
	assert.Error(t, err, "test error")
	assert.Equal(t, false, actual)
}
