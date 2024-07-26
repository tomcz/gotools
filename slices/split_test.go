package slices

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestSplit_Empty(t *testing.T) {
	assert.Nil(t, Split[int](nil, 10))
}

func TestSplit_lessThanPartSize(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	assert.Equal(t, [][]int{in}, Split(in, 10))
}

func TestSplit_equalPartSize(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	assert.Equal(t, [][]int{in}, Split(in, 5))
}

func TestSplit_evenPartSize(t *testing.T) {
	in := []int{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}
	expected := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	assert.Equal(t, expected, Split(in, 3))
}

func TestSplit_oddPartSize(t *testing.T) {
	in := []int{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
		10,
	}
	expected := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	}
	assert.Equal(t, expected, Split(in, 3))
}
