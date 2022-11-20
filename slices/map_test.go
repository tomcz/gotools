package slices

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []uint{1, 2, 3, 4, 5}
	assert.Equal(t, expected, Map(input, func(x int) uint { return uint(x) }))
}

func TestMapErr(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []uint{1, 2, 3, 4, 5}

	actual, err := MapErr(input, func(x int) (uint, error) { return uint(x), nil })
	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}

	actual, err = MapErr(input, func(x int) (uint, error) { return uint(x), errors.New("test error") })
	if assert.EqualError(t, err, "test error") {
		assert.Nil(t, actual)
	}
}
