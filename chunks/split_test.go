package chunks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface_Empty(t *testing.T) {
	assert.Nil(t, SplitInt(nil, 10))
}

func TestInterface_lessThanPartSize(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	assert.Equal(t, [][]int{in}, SplitInt(in, 10))
}

func TestInterface_equalPartSize(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	assert.Equal(t, [][]int{in}, SplitInt(in, 5))
}

func TestInterface_evenPartSize(t *testing.T) {
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
	assert.Equal(t, expected, SplitInt(in, 3))
}

func TestInterface_oddPartSize(t *testing.T) {
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
	assert.Equal(t, expected, SplitInt(in, 3))
}
