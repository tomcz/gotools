package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	expected := []DstType{1, 2, 3, 4}
	actual := ConvertSrcTypeToDstType([]SrcType{1, 2, 3, 4})
	assert.Equal(t, expected, actual)
}

func TestAppend(t *testing.T) {
	expected := []DstType{"foo", 2, 3, 4}
	actual := AppendSrcTypeToDstType([]DstType{"foo"}, 2, 3, 4)
	assert.Equal(t, expected, actual)
}
