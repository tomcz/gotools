package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	expected := []Dst{1, 2, 3, 4}
	actual := ConvertSrcToDst([]Src{1, 2, 3, 4})
	assert.Equal(t, expected, actual)
}

func TestAppend(t *testing.T) {
	expected := []Dst{"foo", 2, 3, 4}
	actual := AppendSrcToDst([]Dst{"foo"}, 2, 3, 4)
	assert.Equal(t, expected, actual)
}
