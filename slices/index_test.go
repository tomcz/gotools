package slices

import (
	"errors"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, 2, Index(src, 3))
	assert.Equal(t, -1, Index(src, 9))
}

func TestIndexOf(t *testing.T) {
	selector := func(v int) bool { return v%2 == 0 }

	src := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, 1, IndexOf(src, selector))

	src = []int{1, 3, 5, 7}
	assert.Equal(t, -1, IndexOf(src, selector))
}

func TestIndexOfErr(t *testing.T) {
	src := []int{1, 2, 3, 4, 5, 6}
	selector := func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err := IndexOfErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, 1, actual)

	src = []int{1, 3, 5, 7}
	selector = func(v int) (bool, error) { return v%2 == 0, nil }
	actual, err = IndexOfErr(src, selector)
	assert.NilError(t, err)
	assert.Equal(t, -1, actual)

	src = []int{1, 2, 3, 4, 5, 6}
	selector = func(v int) (bool, error) { return false, errors.New("test error") }
	actual, err = IndexOfErr(src, selector)
	assert.Error(t, err, "test error")
	assert.Equal(t, -1, actual)
}

func TestIndexBy(t *testing.T) {
	src := []string{
		"one_foo",
		"one_bar",
		"one_wee",
		"two_foo",
	}
	expected := map[string]string{
		"one": "one_wee",
		"two": "two_foo",
	}

	indexFunc := func(value string) string {
		tokens := strings.Split(value, "_")
		return tokens[0]
	}
	assert.DeepEqual(t, expected, IndexBy(src, indexFunc))

	indexFuncOk := func(value string) (string, error) {
		return indexFunc(value), nil
	}
	actual, err := IndexByErr(src, indexFuncOk)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)

	indexFuncBad := func(value string) (string, error) {
		return "", errors.New("test error")
	}
	actual, err = IndexByErr(src, indexFuncBad)
	assert.Error(t, err, "test error")
	assert.Assert(t, is.Nil(actual))
}
