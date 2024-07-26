package slices

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []uint{1, 2, 3, 4, 5}
	mapper := func(x int) uint { return uint(x) }
	assert.Equal(t, expected, Map(input, mapper))
}

func TestMapErr(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []uint{1, 2, 3, 4, 5}

	mapper := func(x int) (uint, error) { return uint(x), nil }
	actual, err := MapErr(input, mapper)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	mapperErr := func(x int) (uint, error) { return uint(x), errors.New("test error") }
	actual, err = MapErr(input, mapperErr)
	assert.Error(t, err, "test error")
	assert.Nil(t, actual)
}
