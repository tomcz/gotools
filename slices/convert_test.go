package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	expected := []interface{}{1, 2, 3, 4}
	actual := ConvertSourceToInterface([]Source{1, 2, 3, 4})
	assert.Equal(t, expected, actual)
}

func TestAppend(t *testing.T) {
	expected := []interface{}{"foo", 2, 3, 4}
	actual := AppendSourceToInterface([]interface{}{"foo"}, 2, 3, 4)
	assert.Equal(t, expected, actual)
}
