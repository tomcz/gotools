package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToInterface(t *testing.T) {
	expected := []interface{}{1, 2, 3, 4}
	actual := ConvertIntToInterface([]int{1, 2, 3, 4})
	assert.Equal(t, expected, actual)
}

func TestAppendToInterface(t *testing.T) {
	expected := []interface{}{"foo", 2, 3, 4}
	actual := AppendIntToInterface([]interface{}{"foo"}, 2, 3, 4)
	assert.Equal(t, expected, actual)
}
