package slices

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []uint{1, 2, 3, 4, 5}
	mapper := func(x int) uint { return uint(x) }
	assert.DeepEqual(t, expected, Map(input, mapper))
}

func TestMapErr(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []uint{1, 2, 3, 4, 5}

	mapper := func(x int) (uint, error) { return uint(x), nil }
	actual, err := MapErr(input, mapper)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)

	mapperErr := func(x int) (uint, error) { return uint(x), errors.New("test error") }
	actual, err = MapErr(input, mapperErr)
	assert.Error(t, err, "test error")
	assert.Assert(t, is.Nil(actual))
}
