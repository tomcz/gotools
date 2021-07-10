package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit_Empty(t *testing.T) {
	assert.Nil(t, SplitSlice(nil, 10))
}

func TestSplit_lessThanPartSize(t *testing.T) {
	in := []Slice{1, 2, 3, 4, 5}
	assert.Equal(t, [][]Slice{in}, SplitSlice(in, 10))
}

func TestSplit_equalPartSize(t *testing.T) {
	in := []Slice{1, 2, 3, 4, 5}
	assert.Equal(t, [][]Slice{in}, SplitSlice(in, 5))
}

func TestSplit_evenPartSize(t *testing.T) {
	in := []Slice{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}
	expected := [][]Slice{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	assert.Equal(t, expected, SplitSlice(in, 3))
}

func TestSplit_oddPartSize(t *testing.T) {
	in := []Slice{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
		10,
	}
	expected := [][]Slice{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	}
	assert.Equal(t, expected, SplitSlice(in, 3))
}
